package synchs

import (
	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
)

func (n *SyncHS) cmdHandler() {
	var suffCmds bool = false
	myID := n.config.GetID()
	pendingCommands := make(map[crypto.Hash]*chain.Command)
	blkSize := n.config.GetBlockSize()
	certMap := make(map[uint64]*msg.BlockCertificate)
	// Setup certificate for the first block
	genesisCert := &msg.BlockCertificate{
		BCert: &msg.Certificate{},
		Data:  &msg.VoteData{},
	}
	genesisCert.Data.View = n.view
	genesisCert.Data.Block = chain.GetGenesis()
	certMap[0] = genesisCert
	for {
		select {
		case cmd, ok := <-n.cmdChannel:
			if !ok {
				log.Error("Command Channel error")
			}
			log.Trace("Handling command:", cmd.String())
			h := cmd.GetHash()
			var exists bool
			// If already present, ignore
			_, exists = pendingCommands[h]
			if !exists {
				pendingCommands[h] = cmd
			} else {
				// We already have this command, continue
				continue
			}
			// If we have sufficient commands then start the blame timer
			suffCmds = uint64(len(pendingCommands)) >= blkSize
			isLeader := n.leader == myID
			log.Trace("Adding command to pending commands buffer")
			if !suffCmds {
				continue
			}
			// Add cmd to pending commands
			// I am not the leader, skip the rest
			if !isLeader {
				// I am not the leader, but there are sufficient commands
				log.Debug("Sufficient commands received. Starting Blame timer")
				go n.startBlameTimer()
				continue
			}
			// As long as we have sufficient commands
			for suffCmds {
				// I am the leader and there are sufficient commands
				n.bc.ChainLock.Lock()
				head := n.bc.Head
				if _, exists = certMap[head]; !exists {
					n.bc.ChainLock.Unlock()
					break
				}
				n.bc.Head++
				n.bc.ChainLock.Unlock()
				// Start building the block
				var cmds []*chain.Command = make([]*chain.Command, blkSize)
				var hashesToRemove []crypto.Hash = make([]crypto.Hash, blkSize)
				var idx = 0
				// First copy the commands onto an array
				for h, cmd := range pendingCommands {
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
					delete(pendingCommands, h)
				}
				// Build a block from pending commands
				nblk := NewBlock(cmds)
				prop := &msg.Proposal{}
				prop.ProposedBlock = nblk
				prop.Cert = certMap[head]
				// Update suffCmds; We may no longer have sufficient commands
				suffCmds = uint64(len(pendingCommands)) >= blkSize
				go n.propose(prop, head+1)
			}
		case cert, ok := <-n.syncObj.certChannel:
			if !ok {
				log.Error("Command Channel error - Sync")
				return
			}
			log.Trace("Received a certificate - Sync")
			certMap[cert.Data.Block.Data.Index] = cert
			if myID != n.leader {
				continue
			}
			if !suffCmds {
				log.Debug("Have a certificate, but have insufficient commands")
				continue
			}
			for suffCmds {
				n.bc.ChainLock.Lock()
				head := n.bc.Head
				if _, exists := certMap[head]; !exists {
					log.Debug("Insufficient certificates. Not Proposing")
					n.bc.ChainLock.Unlock()
					continue
				}
				n.bc.Head++
				n.bc.ChainLock.Unlock()
				// Start building the block
				var cmds []*chain.Command = make([]*chain.Command, blkSize)
				var hashesToRemove []crypto.Hash = make([]crypto.Hash, blkSize)
				var idx = 0
				// First copy the commands onto an array
				for h, cmd := range pendingCommands {
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
					delete(pendingCommands, h)
				}
				// Build a block from pending commands
				nblk := NewBlock(cmds)
				prop := &msg.Proposal{}
				prop.ProposedBlock = nblk
				prop.Cert = certMap[head]
				// Update suffCmds; We may no longer have sufficient commands
				suffCmds = uint64(len(pendingCommands)) >= blkSize
				go n.propose(prop, head+1)
			}
		case prop, ok := <-n.syncObj.propChannel:
			if !ok {
				log.Error("Proposal Channel Error - Sync")
				return
			}
			log.Trace("Received a proposal - Sync")
			if _, exists := certMap[prop.ProposedBlock.Data.Index]; !exists {
				certMap[prop.ProposedBlock.Data.Index] = prop.Cert
			}
		}
	}
}
