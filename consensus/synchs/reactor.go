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
			return
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
			log.Trace("Received a proposal from ", prop.ProposedBlock.Proposer)
			// Send proposal to propose handler
			go n.proposeHandler(prop)
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
			log.Warn("Unspecified msg type", x)
		default:
			log.Warn("Unknown msg type", x)
		}
	}
}
