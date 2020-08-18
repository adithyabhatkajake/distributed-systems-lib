package e2c

import (
	"fmt"

	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	pb "github.com/golang/protobuf/proto"
)

func (e *E2C) react(m []byte) {
	fmt.Println("Received a message of size", len(m))

	inMessage := &msg.E2CMsg{}
	err := pb.Unmarshal(m, inMessage)
	if err != nil {
		panic(err)
	}
	e.msgChannel <- inMessage
}

func (e *E2C) protocol() {
	// Process protocol messages
	for {
		msgIn, ok := <-e.msgChannel
		if !ok {
			fmt.Println("Msg channel error")
		}
		fmt.Println("Received msg", msgIn.String())
		switch x := msgIn.Msg.(type) {
		case *msg.E2CMsg_Cmd:
			fmt.Println("Got a command from client boss!")
			cmd := msgIn.GetCmd()
			fmt.Println("Cmd is:", string(cmd.Cmd))
			// Everyone adds cmd to pending commands
			e.cmdChannel <- cmd
		case *msg.E2CMsg_Prop:
			prop := msgIn.GetProp()
			fmt.Println("Received a propoal from", prop.ProposedBlock.Proposer)
			go e.handleProposal(prop)
		case nil:
			fmt.Println("Unspecified type")
		default:
			fmt.Println("Unknown type", x)
		}
	}
}

func (e *E2C) cmdHandler() {
	for {
		cmd, ok := <-e.cmdChannel
		if !ok {
			fmt.Println("Command Channel error")
		}
		fmt.Println("Handling command:", cmd.String())
		h := cmd.GetHash()
		var exists bool
		e.cmdMutex.Lock()
		_, exists = e.pendingCommands[h]
		if !exists {
			e.pendingCommands[h] = cmd
		}
		e.cmdMutex.Unlock()
		// We already received this command once, skip
		if exists {
			continue
		}
		fmt.Println("Added command to pending commands")
		// I am not the leader, skip the rest
		if e.leader != e.config.GetID() {
			continue
		}
		// If I am the leader, then propose
		go e.propose()
	}
}
