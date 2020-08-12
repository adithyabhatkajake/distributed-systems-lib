package e2c

import (
	"bufio"
	"fmt"

	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	"github.com/libp2p/go-libp2p-core/network"
)

// Implement how to talk to clients
const (
	ClientProtocolID = "e2c/client/0.0.1"
)

// ClientMsgHandler defines how to talk to client messages
func (e *E2C) ClientMsgHandler(s network.Stream) {
	// A buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	// Event Handler
	for {
		// Receive a message from a client and process them
		len, err := rw.Read(buf)
		if err != nil {
			e.errCh <- err
			fmt.Println(err)
			return
		}
		// Send a copy for reacting
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React concurrently
		go e.react(msgBuf)
	}
}
