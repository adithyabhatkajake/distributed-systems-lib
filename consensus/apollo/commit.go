package apollo

import (
	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/apollo"
)

// CommitPrefix takes a block and commits the prefix as per Apollo rule
func (n *Apollo) CommitPrefix(blk *chain.Block) {
	log.Trace("Committing block", blk.String())
	if blk.Data.Index < n.config.GetNumberOfFaultyNodes() {
		return
	}
	startIdx := blk.Data.Index - n.config.GetNumberOfFaultyNodes()
	n.bc.ChainLock.Lock()
	defer n.bc.ChainLock.Unlock()
	for {
		cblk := n.bc.HeightBlockMap[startIdx]
		if cblk.Decision == false {
			cblk.Decision = true
			log.Info("Committing block-", cblk.Data.Index)
			go n.acknowledgeBlockToClients(cblk)
			// Now remove this block from unconfirmed blocks
			delete(n.bc.UnconfirmedBlocks, cblk.GetHashBytes())
			startIdx--
		} else {
			return
		}
		if startIdx == 0 {
			return
		}
	}
}

func (n *Apollo) acknowledgeBlockToClients(blk *chain.Block) {
	log.Debug("Acknowledging clients of commit for block-", blk.Data.Index)
	var err error
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
		climsg := &msg.ApolloMsg{}
		climsg.Msg = &msg.ApolloMsg_Ack{Ack: ack}
		// Tell all the clients, that I have committed this block
		n.ClientBroadcast(climsg)
	}
}
