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
	switch x := inMessage.Msg.(type) {
	case *msg.E2CMsg_Cmd:
		fmt.Println("Got a command from client boss!")
		cmd := inMessage.GetCmd()
		fmt.Println("Cmd is:", string(cmd.Cmd))
		// Add cmd to pending commands
		// If leader, propose
		// Self deliver a proposal
	case nil:
		fmt.Println("Unspecified type")
	default:
		fmt.Println("Unknown type", x)
	}

}

func (e *E2C) protocol() {
	for {
		msg, ok := <-e.msgChannel
		if !ok {
			fmt.Println("Msg channel error")
		}
		fmt.Println("Received msg", msg.String())
	}
}
