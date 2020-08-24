package e2c

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	pb "github.com/golang/protobuf/proto"
)

func (e *E2C) react(m []byte) {
	log.Trace("Received a message of size", len(m))
	inMessage := &msg.E2CMsg{}
	err := pb.Unmarshal(m, inMessage)
	if err != nil {
		log.Error("Received an invalid protocol message", err)
		return
	}
	e.msgChannel <- inMessage
}

func (e *E2C) protocol() {
	// Process protocol messages
	for {
		msgIn, ok := <-e.msgChannel
		if !ok {
			log.Error("Msg channel error")
		}
		log.Trace("Received msg", msgIn.String())
		switch x := msgIn.Msg.(type) {
		case *msg.E2CMsg_Cmd:
			log.Trace("Got a command from client boss!")
			cmd := msgIn.GetCmd()
			log.Trace("Cmd is:", string(cmd.Cmd))
			// Everyone adds cmd to pending commands
			e.cmdChannel <- cmd
		case *msg.E2CMsg_Prop:
			prop := msgIn.GetProp()
			log.Trace("Received a propoal from", prop.ProposedBlock.Proposer)
			go e.handleProposal(prop)
		case *msg.E2CMsg_Npblame:
			blMsg := msgIn.GetNpblame()
			go e.handleNoProgressBlame(blMsg)
		case *msg.E2CMsg_Eqblame:
			_ = msgIn.GetEqblame()
			// TODO
		case nil:
			log.Warn("Unspecified type")
		default:
			log.Warn("Unknown type", x)
		}
	}
}

func (e *E2C) cmdHandler() {
	for {
		cmd, ok := <-e.cmdChannel
		if !ok {
			log.Error("Command Channel error")
		}
		log.Trace("Handling command:", cmd.String())
		h := cmd.GetHash()
		var exists bool
		log.Trace("Trying to acquire cmdMutex lock")
		e.cmdMutex.Lock()
		log.Trace("Acquired cmdMutex lock")
		// If this is the first command, start the blame timer
		log.Trace("Checking if we are adding a command to an empty pendingCommmads buffer")
		if len(e.pendingCommands) == 0 {
			log.Debug("First command received. Starting Blame timer")
			go e.startBlameTimer()
		}
		log.Trace("Adding command to pending commands buffer")
		// Add cmd to pending commands
		_, exists = e.pendingCommands[h]
		if !exists {
			e.pendingCommands[h] = cmd
		}
		e.cmdMutex.Unlock()
		// We already received this command once, skip
		if exists {
			continue
		}
		log.Trace("Added command to pending commands")
		// I am not the leader, skip the rest
		if e.leader != e.config.GetID() {
			continue
		}
		// If I am the leader, then propose
		go e.propose()
	}
}
