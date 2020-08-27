package synchs

import (
	"bytes"
	"time"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
	"github.com/adithyabhatkajake/libe2c/util"
	pb "github.com/golang/protobuf/proto"
)

// Monitor pending commands, if there is any change and the current node is the leader, then start proposing blocks
func (n *SyncHS) propose() {
	log.Trace("Starting a propose step")
	// If there are sufficient commands in pendingCommands, then propose
	n.cmdMutex.Lock()
	blkSize := n.config.GetBlockSize()
	// Insufficient commands to make a block
	if uint64(len(n.pendingCommands)) < blkSize {
		n.cmdMutex.Unlock()
		return
	}
	// We have sufficient blocks by now
	// Do we have a certificate for the previous block?
	n.bc.ChainLock.Lock()
	_, exists := n.certMap[n.bc.Head]
	n.bc.ChainLock.Unlock()
	if !exists {
		n.cmdMutex.Unlock()
		log.Debug("No certificate for head found. Aborting proposal")
		return
	}
	// We have certificate for the head and also sufficient commands
	// Start building the block
	var cmds []*chain.Command = make([]*chain.Command, blkSize)
	var hashesToRemove []crypto.Hash = make([]crypto.Hash, blkSize)
	var idx = 0
	// First copy the commands onto an array
	for h, cmd := range n.pendingCommands {
		cmds[idx] = cmd
		hashesToRemove[idx] = h
		idx++
		// Stop early if we have collected sufficient commands to make a block
		if uint64(idx) == blkSize {
			break
		}
	}
	// Now remove them from the map
	for _, h := range hashesToRemove {
		delete(n.pendingCommands, h)
	}
	// Others can freely mutate the pending commands now, we are done with it
	n.cmdMutex.Unlock()

	prop := &msg.Proposal{}
	blk := NewBlock(cmds)
	blk.Proposer = n.config.GetID()
	// We are mutating the chain, acquire the lock first
	{
		n.bc.ChainLock.Lock()
		n.certMapLock.RLock()
		// Set index
		blk.Data.Index = n.bc.Head + 1
		// Set previous hash to the current head
		blk.Data.PrevHash = n.bc.HeightBlockMap[n.bc.Head].GetBlockHash()
		prop.Cert = n.certMap[n.bc.Head]
		// Increment chain head for future proposals
		n.bc.Head++
		// Set old head, so that parallel proposals have the correct prevHash
		n.bc.HeightBlockMap[n.bc.Head] = blk
		// Set Hash
		blk.BlockHash = blk.GetHash().GetBytes()
		// Set unconfirmed Blocks
		n.bc.UnconfirmedBlocks[blk.GetHashBytes()] = blk
		// e.bc.HeightBlockMap[e.bc.head] = blk
		// The chain is free for mutation from others now
		n.certMapLock.RUnlock()
		n.bc.ChainLock.Unlock()
	}
	// Sign
	data, err := pb.Marshal(blk.Data)
	if err != nil {
		log.Error("Error in marshalling block data during proposal")
		panic(err)
	}
	sig, err := n.config.GetMyKey().Sign(data)
	if err != nil {
		log.Error("Error in signing a block during proposal")
		panic(err)
	}
	blk.Signature = sig
	prop.View = n.view
	prop.ProposedBlock = blk
	log.Trace("Finished Proposing")
	// Ship proposal to processing

	relayMsg := &msg.SyncHSMsg{}
	relayMsg.Msg = &msg.SyncHSMsg_Prop{Prop: prop}
	n.Broadcast(relayMsg) // Leader sends new block to all the other nodes
	// Leader should also vote
	n.voteForBlock(blk)
	// Start 2\delta timer
	n.startBlockTimer(prop.ProposedBlock)
}

// Deal with the proposal
func (n *SyncHS) handleProposal(prop *msg.Proposal) {
	log.Trace("Handling proposal", prop.ProposedBlock.Data.Index)
	if !prop.ProposedBlock.IsValid() {
		log.Warn("Invalid block")
		return
	}
	data, err := pb.Marshal(prop.ProposedBlock.Data)
	if err != nil {
		log.Error("Proposal error:", err)
		return
	}
	correct, err := n.config.GetPubKeyFromID(n.leader).Verify(data,
		prop.ProposedBlock.Signature)
	if !correct {
		log.Error("Incorrect signature for proposal", prop)
		return
	}
	// Check block certificate for non-genesis blocks
	if !n.IsCertValid(prop.Cert) {
		log.Error("Invalid certificate received for block", prop.ProposedBlock.Data.Index)
		return
	}
	var blk *chain.Block
	var exists bool
	{
		// First check for equivocation
		n.bc.ChainLock.RLock()
		blk, exists = n.bc.HeightBlockMap[prop.ProposedBlock.Data.Index]
		n.bc.ChainLock.RUnlock()
	}
	if exists &&
		!bytes.Equal(prop.ProposedBlock.GetBlockHash(), blk.GetBlockHash()) {
		// Equivocation
		log.Warn("Equivocation detected.", blk, prop.ProposedBlock)
		// TODO trigger view change
		return
	}
	if exists {
		// Duplicate block received,
		// we have already committed this block, IGNORE
		return
	}
	{
		n.bc.ChainLock.RLock()
		_, exists = n.bc.UnconfirmedBlocks[prop.ProposedBlock.GetHashBytes()]
		n.bc.ChainLock.RUnlock()
	}
	if exists {
		// Duplicate block received,
		// We have already received this proposal, IGNORE
		return
	}
	n.addNewBlock(prop.ProposedBlock)
	n.ensureBlockIsDelivered(prop.ProposedBlock)

	// Remove cmds proposed from pending commands
	n.cmdMutex.Lock()
	for _, cmd := range prop.ProposedBlock.Data.Cmds {
		h := cmd.GetHash()
		delete(n.pendingCommands, h)
	}
	n.cmdMutex.Unlock()
	// Vote for the proposal
	n.voteForBlock(prop.ProposedBlock)
	// Start 2\delta timer
	n.startBlockTimer(prop.ProposedBlock)
	// Stop blame timer, since we got a valid proposal
	// During commit, if pending commands is empty, we will restart the blame timer
	go n.stopBlameTimer()
}

func (n *SyncHS) voteForBlock(blk *chain.Block) {
	v := &msg.Vote{}
	v.Origin = n.config.GetID()
	v.Data = &msg.VoteData{}
	v.Data.Block = blk
	v.Data.View = n.view
	data, err := pb.Marshal(v.Data)
	if err != nil {
		log.Error("Error marshing vote data during voting")
		log.Error(err)
		return
	}
	v.Signature, err = n.config.GetMyKey().Sign(data)
	if err != nil {
		log.Error("Error signing vote")
		log.Error(err)
		return
	}
	voteMsg := &msg.SyncHSMsg{}
	voteMsg.Msg = &msg.SyncHSMsg_Vote{Vote: v}
	n.Broadcast(voteMsg) // Send vote to all the nodes
	// Handle my own vote
	go func() {
		n.voteChannel <- v
	}()
}

// NewBlock creates a new block from the commands received.
func NewBlock(cmds []*chain.Command) *chain.Block {
	b := &chain.Block{}
	b.Data = &chain.BlockData{}
	b.Data.Cmds = cmds
	b.Decision = false
	return b
}

func (n *SyncHS) ensureBlockIsDelivered(blk *chain.Block) {
	var exists bool
	var parentblk *chain.Block
	// Ensure that all the parents are delivered first.
	parentIdx := blk.Data.Index - 1
	// Wait for parents to be delivered first
	for tries := 30; tries > 0; tries-- {
		<-time.After(time.Second)
		n.bc.ChainLock.RLock()
		parentblk, exists = n.bc.HeightBlockMap[parentIdx]
		n.bc.ChainLock.RUnlock()
		if exists &&
			!bytes.Equal(parentblk.BlockHash, blk.Data.PrevHash) {
			// This block is delivered.
			log.Warn("Block  ", blk.Data.Index, " extending wrong parent.\n",
				"Wanted Parent Block:", util.BytesToHexString(parentblk.BlockHash),
				"Found Parent Block:", util.BytesToHexString(blk.Data.PrevHash))
			return
		}
		if exists {
			// The parent of the proposed block is the same as the block we have at the parent's position, CONTINUE
			break
		}
	}
	if !exists {
		// The parents are not delivered, so we cant process this block
		// Return
		log.Warn("Parents not delivered, aborting this proposal")
		return
	}
	// All parents are delivered, lets break out!!
	log.Trace("All parents are delivered")
}

func (n *SyncHS) startBlockTimer(blk *chain.Block) {
	var err error
	// Start 2delta timer
	timer := util.NewTimer(func() {
		log.Info("Committing block-", blk.Data.Index)
		// We have committed this block
		blk.Decision = true
		// Let the client know that we committed this block
		for _, cmd := range blk.Data.Cmds {
			ack := &msg.CommitAck{}
			cmdHash := cmd.GetHash()
			ack.CmdHash = cmdHash.GetBytes()
			ack.Id = n.config.GetID()
			ack.Signature, err = n.config.GetMyKey().Sign(ack.CmdHash)
			log.Trace("Sending ack ", ack.CmdHash, " to clients")
			if err != nil {
				log.Error("Error sending ack ", ack.CmdHash, " to clients")
				continue
			}
			synchsmsg := &msg.SyncHSMsg{}
			synchsmsg.Msg = &msg.SyncHSMsg_Ack{Ack: ack}
			// Tell all the clients, that I have committed this block
			n.ClientBroadcast(synchsmsg)
			// Now remove this block from unconfirmed blocks
			n.bc.ChainLock.Lock()
			delete(n.bc.UnconfirmedBlocks, blk.GetHashBytes())
			n.bc.ChainLock.Unlock()
		}
	})
	log.Info("Started timer for block-", blk.Data.Index)
	timer.SetTime(n.config.GetCommitWaitTime())
	n.addNewTimer(blk.Data.Index, timer)
	timer.Start()
}

func (n *SyncHS) addNewBlock(blk *chain.Block) {
	// Otherwise, add the current block to map
	n.bc.ChainLock.Lock()
	n.bc.HeightBlockMap[blk.Data.Index] = blk
	n.bc.UnconfirmedBlocks[blk.GetHashBytes()] =
		blk
	n.bc.ChainLock.Unlock()
}

func (n *SyncHS) addNewTimer(pos uint64, timer *util.Timer) {
	n.timerLock.Lock()
	n.timerMaps[pos] = timer
	n.timerLock.Unlock()
}
