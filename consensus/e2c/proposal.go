package e2c

import (
	"bytes"
	"fmt"
	"time"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/crypto"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	"github.com/adithyabhatkajake/libe2c/util"
	pb "github.com/golang/protobuf/proto"
)

// Monitor pending commands, if there is any change and the current node is the leader, then start proposing blocks
func (e *E2C) propose() {
	fmt.Println("Starting a propose step")
	// If there are sufficient commands in pendingCommands, then propose
	e.cmdMutex.Lock()
	blkSize := e.config.GetBlockSize()
	// Insufficient commands to make a block
	if uint64(len(e.pendingCommands)) < blkSize {
		e.cmdMutex.Unlock()
		return
	}
	// We have sufficient blocks, let us propose now
	var cmds []*chain.Command = make([]*chain.Command, blkSize)
	var hashesToRemove []crypto.Hash = make([]crypto.Hash, blkSize)
	var idx = 0
	// First copy the commands onto an array
	for h, cmd := range e.pendingCommands {
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
		e.cmdMutex.Unlock()
		return
	}
	// Now remove them from the map
	for _, h := range hashesToRemove {
		delete(e.pendingCommands, h)
	}
	// Others can freely mutate the pending commands now, we are done with it
	e.cmdMutex.Unlock()

	blk := NewBlock(cmds)
	blk.Proposer = e.config.GetID()
	// We are mutating the chain, acquire the lock first
	e.bc.ChainLock.Lock()
	// Set index
	blk.Data.Index = e.bc.Head + 1
	// Set previous hash to the current head
	blk.Data.PrevHash = e.bc.HeightBlockMap[e.bc.Head].GetBlockHash()
	// Increment chain head for future proposals
	e.bc.Head++
	// Set Hash
	blk.BlockHash = blk.GetHash().GetBytes()
	// After setting hash, update block head
	e.bc.HeightBlockMap[e.bc.Head] = blk
	// The chain is free for mutation from others now
	e.bc.ChainLock.Unlock()
	// Sign
	data, err := pb.Marshal(blk.Data)
	if err != nil {
		fmt.Println("Error in marshalling block data during proposal")
		panic(err)
	}
	sig, err := e.config.PvtKey.Sign(data)
	if err != nil {
		fmt.Println("Error in signing a block during proposal")
		panic(err)
	}
	blk.Signature = sig
	// Send proposal to all other nodes
	sendMsg := &msg.E2CMsg{}
	prop := &msg.Proposal{}
	prop.ProposedBlock = blk
	sendMsg.Msg = &msg.E2CMsg_Prop{Prop: prop}
	e.Broadcast(sendMsg)
	fmt.Println("Finished Proposing")
}

func (e *E2C) handleProposal(prop *msg.Proposal) {
	fmt.Println("Handling proposal", prop.ProposedBlock.Data.Index)
	if !prop.ProposedBlock.IsValid() {
		fmt.Println("Invalid block")
		return
	}
	data, err := pb.Marshal(prop.ProposedBlock.Data)
	if err != nil {
		fmt.Println("Proposal error:", err)
		return
	}
	correct, err := e.config.GetPubKeyFromID(e.leader).Verify(data,
		prop.ProposedBlock.Signature)
	if !correct {
		fmt.Println("Incorrect signature for proposal")
		return
	}
	// First check for equivocation
	e.bc.ChainLock.Lock()
	blk, exists := e.bc.HeightBlockMap[prop.ProposedBlock.Data.Index]
	e.bc.ChainLock.Unlock()
	if exists &&
		!bytes.Equal(prop.ProposedBlock.GetBlockHash(), blk.GetBlockHash()) {
		// Equivocation
		fmt.Println("Equivocation detected.", blk, prop.ProposedBlock)
		// TODO trigger view change
		return
	}
	if exists {
		// Duplicate block received,
		// we have already committed this block, IGNORE
		return
	}
	e.bc.ChainLock.Lock()
	_, exists = e.bc.UnconfirmedBlocks[prop.ProposedBlock.GetHashBytes()]
	e.bc.ChainLock.Unlock()
	if exists {
		// Duplicate block received,
		// We have already received this proposal, IGNORE
		return
	}
	// Otherwise, add the current block to map
	e.bc.ChainLock.Lock()
	e.bc.HeightBlockMap[prop.ProposedBlock.Data.Index] = prop.ProposedBlock
	e.bc.UnconfirmedBlocks[prop.ProposedBlock.GetHashBytes()] =
		prop.ProposedBlock
	e.bc.ChainLock.Unlock()
	// Ensure that all the parents are delivered first.
	parentIdx := prop.ProposedBlock.Data.Index - 1
	// Wait for parents to be delivered first
	for tries := 30; tries > 0; tries-- {
		<-time.After(time.Second)
		e.bc.ChainLock.Lock()
		blk, exists = e.bc.HeightBlockMap[parentIdx]
		e.bc.ChainLock.Unlock()
		if exists &&
			!bytes.Equal(blk.BlockHash, prop.ProposedBlock.Data.PrevHash) {
			// This block is delivered.
			fmt.Println("Block extending wrong parent.")
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
		fmt.Println("Parents not delivered, aborting this proposal")
		return
	}
	// All parents are delivered, lets break out!!
	fmt.Println("All parents are delivered")

	// Start 2delta timer
	timer := util.NewTimer(func() {
		fmt.Println("Committing block", prop.ProposedBlock)
	})
	timer.SetTime(e.config.GetCommitWaitTime())
	e.timerMaps[prop.ProposedBlock.Data.Index] = timer
	timer.Start()
}

// NewBlock creates a new block from the commands received.
func NewBlock(cmds []*chain.Command) *chain.Block {
	b := &chain.Block{}
	b.Data = &chain.BlockData{}
	b.Data.Cmds = cmds
	b.Decision = false
	return b
}
