package apollo

import (
	"bytes"
	"time"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/apollo"
	"github.com/adithyabhatkajake/libe2c/util"
	pb "github.com/golang/protobuf/proto"
)

// Monitor pending commands, if there is any change and the current node is the leader, then start proposing blocks
func (n *Apollo) propose() {
	log.Trace("Starting a propose step")
	n.leaderLock.Lock() // Get a writing lock
	defer n.leaderLock.Unlock()
	if n.leader != n.config.GetID() {
		return
	}
	// If there are sufficient commands in pendingCommands, then propose
	n.cmdMutex.Lock()
	blkSize := n.config.GetBlockSize()
	// Insufficient commands to make a block
	if uint64(len(n.pendingCommands)) < blkSize {
		n.cmdMutex.Unlock()
		return
	}
	// We have sufficient blocks, let us propose now
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
	// There were insufficient pending commands, wait for more commands
	if uint64(idx) < blkSize {
		n.cmdMutex.Unlock()
		return
	}
	// Now remove them from the map
	for _, h := range hashesToRemove {
		delete(n.pendingCommands, h)
	}
	// Others can freely mutate the pending commands now, we are done with it
	n.cmdMutex.Unlock()

	blk := NewBlock(cmds)
	blk.Proposer = n.config.GetID()
	// We are mutating the chain, acquire the lock first
	{
		n.bc.ChainLock.Lock()
		// Set index
		blk.Data.Index = n.bc.Head + 1
		// Set previous hash to the current head
		blk.Data.PrevHash = n.bc.HeightBlockMap[n.bc.Head].GetBlockHash()
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
		n.bc.ChainLock.Unlock()
	}
	// Sign
	data, err := pb.Marshal(blk.Data)
	if err != nil {
		log.Error("Error in marshalling block data during proposal")
		panic(err)
	}
	sig, err := n.config.PvtKey.Sign(data)
	if err != nil {
		log.Error("Error in signing a block during proposal")
		panic(err)
	}
	blk.Signature = sig
	prop := &msg.Proposal{}
	prop.ProposedBlock = blk
	log.Trace("Finished Proposing")
	// Ship proposal to processing

	relayMsg := &msg.ApolloMsg{}
	relayMsg.Msg = &msg.ApolloMsg_Prop{Prop: prop}
	n.Broadcast(relayMsg) // Leader sends new block to all the other nodes
	// Start 2\delta timer
	n.leader = (n.leader + 1) % n.config.GetNumNodes() // Change leader, so I don't propose again
}

// Deal with the proposal
// This will also relay the proposal to all other nodes
func (n *Apollo) handleProposal(prop *msg.Proposal) {
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
	// Signature check
	sender := n.senderForBlockIndex(prop.ProposedBlock.Data.Index)
	correct, err := n.config.GetPubKeyFromID(sender).Verify(data,
		prop.ProposedBlock.Signature)
	if !correct {
		log.Error("Incorrect signature for proposal")
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
		if n.bc.Head < prop.ProposedBlock.Data.Index {
			n.bc.Head = prop.ProposedBlock.Data.Index
		}
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
	relayMsg := &msg.ApolloMsg{}
	relayMsg.Msg = &msg.ApolloMsg_Prop{Prop: prop}
	n.changeLeader() // Change the leader, since we received a block for this round
	nextl := n.currentLeader()
	if n.config.GetID() == nextl {
		// If I am the next leader, then propose the next block
		go n.propose()
	} else {
		// I am not the next leader, send the block to the next leader
		n.Send(relayMsg, nextl) // Relay it to all the other nodes
	}

	log.Trace("Finished relaying the proposal to the next leader")
	// TODO - Stop blame timer for this leader, since we got a valid proposal
	// go n.stopBlameTimer()
	// TODO - Commit prefix
	n.CommitPrefix(prop.ProposedBlock)
}

// NewBlock creates a new block from the commands received.
func NewBlock(cmds []*chain.Command) *chain.Block {
	b := &chain.Block{}
	b.Data = &chain.BlockData{}
	b.Data.Cmds = cmds
	b.Decision = false
	return b
}

func (n *Apollo) ensureBlockIsDelivered(blk *chain.Block) {
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

func (n *Apollo) addNewBlock(blk *chain.Block) {
	// Otherwise, add the current block to map
	n.bc.ChainLock.Lock()
	n.bc.HeightBlockMap[blk.Data.Index] = blk
	n.bc.UnconfirmedBlocks[blk.GetHashBytes()] =
		blk
	n.bc.ChainLock.Unlock()
}
