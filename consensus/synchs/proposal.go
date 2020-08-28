package synchs

import (
	"bytes"
	"time"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
	"github.com/adithyabhatkajake/libe2c/util"
	pb "github.com/golang/protobuf/proto"
)

// Monitor pending commands, if there is any change and the current node is the leader, then start proposing blocks
func (n *SyncHS) propose(prop *msg.Proposal, pIdx uint64) {
	log.Trace("Starting a propose step")
	prop.ProposedBlock.Proposer = n.config.GetID()
	prop.ProposedBlock.Data.Index = pIdx
	// We are mutating the chain, acquire the lock first
	{
		n.bc.ChainLock.Lock()
		// Set previous hash to the current head
		prop.ProposedBlock.Data.PrevHash = n.bc.HeightBlockMap[pIdx-1].GetBlockHash()
		// Set old head, so that parallel proposals have the correct prevHash
		n.bc.HeightBlockMap[pIdx] = prop.ProposedBlock
		// Set Hash
		prop.ProposedBlock.BlockHash = prop.ProposedBlock.GetHash().GetBytes()
		// Set unconfirmed Blocks
		n.bc.UnconfirmedBlocks[prop.ProposedBlock.GetHashBytes()] = prop.ProposedBlock
		// e.bc.HeightBlockMap[e.bc.head] = blk
		// The chain is free for mutation from others now
		n.bc.ChainLock.Unlock()
	}
	// Sign
	data, err := pb.Marshal(prop.ProposedBlock.Data)
	if err != nil {
		log.Error("Error in marshalling block data during proposal")
		panic(err)
	}
	sig, err := n.config.GetMyKey().Sign(data)
	if err != nil {
		log.Error("Error in signing a block during proposal")
		panic(err)
	}
	prop.ProposedBlock.Signature = sig
	prop.View = n.view
	log.Trace("Finished Proposing")
	// Ship proposal to processing

	relayMsg := &msg.SyncHSMsg{}
	relayMsg.Msg = &msg.SyncHSMsg_Prop{Prop: prop}
	// Leader sends new block to all the other nodes
	go n.Broadcast(relayMsg)
	// Leader should also vote
	go n.voteForBlock(prop.ProposedBlock)
	// Start 2\delta timer
	go n.startBlockTimer(prop.ProposedBlock)
}

// Deal with the proposal
func (n *SyncHS) proposeHandler(prop *msg.Proposal) {
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
	// Send the proposal to syncer, so that it can remove the candidate blocks and proposals
	go func() {
		n.syncObj.propChannel <- prop
	}()
	// Vote for the proposal
	go n.voteForBlock(prop.ProposedBlock)
	// Start 2\delta timer
	go n.startBlockTimer(prop.ProposedBlock)
	// Stop blame timer, since we got a valid proposal
	// During commit, if pending commands is empty, we will restart the blame timer
	go n.stopBlameTimer()

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
