package e2c

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p"

	config "github.com/adithyabhatkajake/libe2c/config/e2c"
	"github.com/adithyabhatkajake/libe2c/net"

	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"

	"github.com/libp2p/go-libp2p-core/network"
)

const (
	// ProtocolID is the ID for E2C Protocol
	ProtocolID = "e2c/e2c/0.0.1"
)

// Init implements the Protocol interface
func (e *E2C) Init(c *config.NodeConfig) {
	e.config = c
	e.leader = DefaultLeaderID
}

// Setup sets up the network components
func (e *E2C) Setup(n *net.Network) error {
	e.host = n.H
	e.ctx = n.Ctx
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(e.config.GetClientListenAddr()),
		libp2p.Identity(e.config.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}
	e.cliHost = host
	e.pMap = n.PeerMap
	e.msgChannel = make(chan msg.E2CMsg, 100)
	e.streamMap = make(map[uint64]network.Stream)
	// How to react to Protocol Messages
	e.host.SetStreamHandler(ProtocolID, e.ProtoMsgHandler)

	// How to react to Client Messages
	e.cliHost.SetStreamHandler(ClientProtocolID, e.ClientMsgHandler)

	// Connect to all the other nodes talking E2C protocol
	for idx, p := range e.pMap {
		fmt.Println("Attempting to open a stream with", p, "using protocol", ProtocolID)
		retries := 30
		for i := retries; i > 0; i-- {
			s, err := e.host.NewStream(e.ctx, p.ID, ProtocolID)
			if err != nil {
				fmt.Println("Error connecting to peers:", err)
				fmt.Println("Retrying in a minute")
				<-time.After(time.Minute)
				continue
			}
			e.streamMap[idx] = s
			fmt.Println("Connected to Node-#", idx)
			break
		}
	}
	fmt.Println("Setup Finished. Ready to do SMR:)")

	return nil
}

// Start implements the Protocol Interface
func (e *E2C) Start() {
	// Start E2C Protocol
	e.protocol()
	// Run the protocol for 10 minutes before shutting down
	<-time.After(10 * time.Minute)
}

// ProtoMsgHandler reacts to all protocol messages in the network
func (e *E2C) ProtoMsgHandler(s network.Stream) {
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
