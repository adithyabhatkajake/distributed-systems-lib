package synchs

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
	pb "github.com/golang/protobuf/proto"
)

func (n *SyncHS) react(m []byte) {
	log.Trace("Received a message of size", len(m))
	inMessage := &msg.SyncHSMsg{}
	err := pb.Unmarshal(m, inMessage)
	if err != nil {
		log.Error("Received an invalid protocol message", err)
		return
	}
	n.msgChannel <- inMessage
}

func (n *SyncHS) protocol() {
	// Process protocol messages
	for {
		msgIn, ok := <-n.msgChannel
		if !ok {
			log.Error("Msg channel error")
		}
		log.Trace("Received msg", msgIn.String())
		switch x := msgIn.Msg.(type) {
		case *msg.SyncHSMsg_Cmd:
			log.Trace("Got a command from client boss!")
			cmd := msgIn.GetCmd()
			log.Trace("Cmd is:", string(cmd.Cmd))
			// Everyone adds cmd to pending commands
			n.cmdChannel <- cmd
		case *msg.SyncHSMsg_Prop:
			prop := msgIn.GetProp()
			log.Trace("Received a propoal from", prop.ProposedBlock.Proposer)
			go n.handleProposal(prop)
		case *msg.SyncHSMsg_Npblame:
			blMsg := msgIn.GetNpblame()
			go n.handleNoProgressBlame(blMsg)
		case *msg.SyncHSMsg_Eqblame:
			_ = msgIn.GetEqblame()
			// TODO
		case *msg.SyncHSMsg_Vote:
			vote := msgIn.GetVote()
			n.voteChannel <- vote
		case nil:
			log.Warn("Unspecified type")
		default:
			log.Warn("Unknown type", x)
		}
	}
}

func (n *SyncHS) cmdHandler() {
	for {
		cmd, ok := <-n.cmdChannel
		if !ok {
			log.Error("Command Channel error")
		}
		log.Trace("Handling command:", cmd.String())
		h := cmd.GetHash()
		var exists bool
		log.Trace("Trying to acquire cmdMutex lock")
		n.cmdMutex.Lock()
		log.Trace("Acquired cmdMutex lock")
		// If this is the first command, start the blame timer
		log.Trace("Checking if we are adding a command to an empty pendingCommmads buffer")
		if len(n.pendingCommands) == 0 {
			log.Debug("First command received. Starting Blame timer")
			go n.startBlameTimer()
		}
		log.Trace("Adding command to pending commands buffer")
		// Add cmd to pending commands
		_, exists = n.pendingCommands[h]
		if !exists {
			n.pendingCommands[h] = cmd
		}
		n.cmdMutex.Unlock()
		// We already received this command once, skip
		if exists {
			continue
		}
		log.Trace("Added command to pending commands")
		// I am not the leader, skip the rest
		if n.leader != n.config.GetID() {
			continue
		}
		// If I am the leader, then propose
		go n.propose()
	}
}
