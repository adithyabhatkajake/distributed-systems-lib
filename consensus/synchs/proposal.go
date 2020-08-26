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
func (shs *SyncHS) propose() {
	log.Trace("Starting a propose step")
	// If there are sufficient commands in pendingCommands, then propose
	shs.cmdMutex.Lock()
	blkSize := shs.config.GetBlockSize()
	// Insufficient commands to make a block
	if uint64(len(shs.pendingCommands)) < blkSize {
		shs.cmdMutex.Unlock()
		return
	}
	// We have sufficient blocks, let us propose now
	var cmds []*chain.Command = make([]*chain.Command, blkSize)
	var hashesToRemove []crypto.Hash = make([]crypto.Hash, blkSize)
	var idx = 0
	// First copy the commands onto an array
	for h, cmd := range shs.pendingCommands {
		cmds[idx] = cmd
		hashesToRemove[idx] = h
		idx++
		// Stop early if we have collected sufficient commands to make a block
		if uint64(idx) == blkSize {
			break
		}
	}
	// There were insufficient pending commands, wait for more commands
	if uint64(idx) < blkSize {
		shs.cmdMutex.Unlock()
		return
	}
	// Now remove them from the map
	for _, h := range hashesToRemove {
		delete(shs.pendingCommands, h)
	}
	// Others can freely mutate the pending commands now, we are done with it
	shs.cmdMutex.Unlock()

	blk := NewBlock(cmds)
	blk.Proposer = shs.config.GetID()
	// We are mutating the chain, acquire the lock first
	{
		shs.bc.ChainLock.Lock()
		// Set index
		blk.Data.Index = shs.bc.Head + 1
		// Set previous hash to the current head
		blk.Data.PrevHash = shs.bc.HeightBlockMap[shs.bc.Head].GetBlockHash()
		// Increment chain head for future proposals
		shs.bc.Head++
		// Set old head, so that parallel proposals have the correct prevHash
		shs.bc.HeightBlockMap[shs.bc.Head] = blk
		// Set Hash
		blk.BlockHash = blk.GetHash().GetBytes()
		// Set unconfirmed Blocks
		shs.bc.UnconfirmedBlocks[blk.GetHashBytes()] = blk
		// e.bc.HeightBlockMap[e.bc.head] = blk
		// The chain is free for mutation from others now
		shs.bc.ChainLock.Unlock()
	}
	// Sign
	data, err := pb.Marshal(blk.Data)
	if err != nil {
		log.Error("Error in marshalling block data during proposal")
		panic(err)
	}
	sig, err := shs.config.PvtKey.Sign(data)
	if err != nil {
		log.Error("Error in signing a block during proposal")
		panic(err)
	}
	blk.Signature = sig
	prop := &msg.Proposal{}
	prop.ProposedBlock = blk
	log.Trace("Finished Proposing")
	// Ship proposal to processing

	relayMsg := &msg.SyncHSMsg{}
	relayMsg.Msg = &msg.SyncHSMsg_Prop{Prop: prop}
	shs.Broadcast(relayMsg) // Leader sends new block to all the other nodes
	// Start 2\delta timer
	shs.startBlockTimer(prop.ProposedBlock)
}

// Deal with the proposal
// This will also relay the proposal to all other nodes
func (shs *SyncHS) handleProposal(prop *msg.Proposal) {
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
	correct, err := shs.config.GetPubKeyFromID(shs.leader).Verify(data,
		prop.ProposedBlock.Signature)
	if !correct {
		log.Error("Incorrect signature for proposal")
		return
	}
	var blk *chain.Block
	var exists bool
	{
		// First check for equivocation
		shs.bc.ChainLock.RLock()
		blk, exists = shs.bc.HeightBlockMap[prop.ProposedBlock.Data.Index]
		shs.bc.ChainLock.RUnlock()
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
		shs.bc.ChainLock.RLock()
		_, exists = shs.bc.UnconfirmedBlocks[prop.ProposedBlock.GetHashBytes()]
		shs.bc.ChainLock.RUnlock()
	}
	if exists {
		// Duplicate block received,
		// We have already received this proposal, IGNORE
		return
	}
	shs.addNewBlock(prop.ProposedBlock)
	shs.ensureBlockIsDelivered(prop.ProposedBlock)

	// Remove cmds proposed from pending commands
	shs.cmdMutex.Lock()
	for _, cmd := range prop.ProposedBlock.Data.Cmds {
		h := cmd.GetHash()
		delete(shs.pendingCommands, h)
	}
	shs.cmdMutex.Unlock()
	relayMsg := &msg.SyncHSMsg{}
	relayMsg.Msg = &msg.SyncHSMsg_Prop{Prop: prop}
	shs.Broadcast(relayMsg) // Relay it to all the other nodes
	log.Trace("Finished relaying the proposal")
	// Start 2\delta timer
	shs.startBlockTimer(prop.ProposedBlock)
	// Stop blame timer, since we got a valid proposal
	// During commit, if pending commands is empty, we will restart the blame timer
	go shs.stopBlameTimer()
}

// NewBlock creates a new block from the commands received.
func NewBlock(cmds []*chain.Command) *chain.Block {
	b := &chain.Block{}
	b.Data = &chain.BlockData{}
	b.Data.Cmds = cmds
	b.Decision = false
	return b
}

func (shs *SyncHS) ensureBlockIsDelivered(blk *chain.Block) {
	var exists bool
	var parentblk *chain.Block
	// Ensure that all the parents are delivered first.
	parentIdx := blk.Data.Index - 1
	// Wait for parents to be delivered first
	for tries := 30; tries > 0; tries-- {
		<-time.After(time.Second)
		shs.bc.ChainLock.RLock()
		parentblk, exists = shs.bc.HeightBlockMap[parentIdx]
		shs.bc.ChainLock.RUnlock()
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

func (shs *SyncHS) startBlockTimer(blk *chain.Block) {
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
			ack.Id = shs.config.GetID()
			ack.Signature, err = shs.config.GetMyKey().Sign(ack.CmdHash)
			log.Trace("Sending ack ", ack.CmdHash, " to clients")
			if err != nil {
				log.Error("Error sending ack ", ack.CmdHash, " to clients")
				continue
			}
			synchsmsg := &msg.SyncHSMsg{}
			synchsmsg.Msg = &msg.SyncHSMsg_Ack{Ack: ack}
			// Tell all the clients, that I have committed this block
			shs.ClientBroadcast(synchsmsg)
			// Now remove this block from unconfirmed blocks
			shs.bc.ChainLock.Lock()
			delete(shs.bc.UnconfirmedBlocks, blk.GetHashBytes())
			shs.bc.ChainLock.Unlock()
		}
	})
	log.Info("Started timer for block-", blk.Data.Index)
	timer.SetTime(shs.config.GetCommitWaitTime())
	shs.addNewTimer(blk.Data.Index, timer)
	timer.Start()
}

func (shs *SyncHS) addNewBlock(blk *chain.Block) {
	// Otherwise, add the current block to map
	shs.bc.ChainLock.Lock()
	shs.bc.HeightBlockMap[blk.Data.Index] = blk
	shs.bc.UnconfirmedBlocks[blk.GetHashBytes()] =
		blk
	shs.bc.ChainLock.Unlock()
}

func (shs *SyncHS) addNewTimer(pos uint64, timer *util.Timer) {
	shs.timerLock.Lock()
	shs.timerMaps[pos] = timer
	shs.timerLock.Unlock()
}
