package e2c

import (
	"fmt"
	"io"

	"github.com/adithyabhatkajake/libe2c/config"
	"github.com/adithyabhatkajake/libe2c/consensus/e2c/msg"
	"github.com/adithyabhatkajake/libe2c/net"

	"github.com/libp2p/go-libp2p-core/network"
)

const (
	// ID of this message type
	ID = "e2c/e2c/0.0.1"
)

// Init implements the Protocol interface
func (e *E2C) Init(c *config.NodeConfig) {
	e.protConfig = c.Config.ProtConfig
	e.leader = DefaultLeaderID
}

// Setup sets up the network components
func (e *E2C) Setup(n *net.Network) error {
	e.host = n.H
	e.ctx = n.Ctx
	e.pMap = n.PeerMap
	e.host.SetStreamHandler(ID, e.MsgHandler)
	for idx, p := range e.pMap {
		s, err := e.host.NewStream(e.ctx, p.ID, ID)
		if err != nil {
			return err
		}
		e.streamMap[idx] = s
		fmt.Println("Connected to Node-#", idx)
	}
	return nil
}

// Start implements the Protocol Interface
func (e *E2C) Start() {

}

// MsgHandler handles all messages to/from the network
func (e *E2C) MsgHandler(s network.Stream) {
	e.errCh = make(chan error, 1)
	defer close(e.errCh)

	// Run a parallel goroutine that closes the channel if any error is detected
	go func() {
		select {
		case err, ok := <-e.errCh:
			if ok {
				fmt.Println("Type-I", err, ok)
			} else {
				fmt.Println("Type-II", err, ok)
			}
		}
		s.Reset()
	}()

	// A global buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	// Event Handler
	for {
		// Receive a message from anyone and process them
		len, err := io.ReadFull(s, buf)
		if err != nil {
			e.errCh <- err
			return
		}
		// Use a copy of the message and send it to off for processing
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React to the message in parallel and continue
		go e.react(msgBuf)
	}
}
