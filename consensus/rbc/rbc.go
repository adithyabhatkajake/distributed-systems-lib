package rbc

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/adithyabhatkajake/libe2c/config"
	"github.com/adithyabhatkajake/libe2c/net"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

/* Reliable Broadcast Protocol Implementation */

const (
	// ID of this protocol
	ID = "e2c/rbc/0.0.1"
)

var (
	// SenderID is the sender for the reliable broadcast
	SenderID uint64 = 1
)

// RBC implements Consensus Protocol
type RBC struct {
	host     host.Host
	ctx      context.Context
	pMap     map[uint64]*peerstore.AddrInfo
	isSender bool
}

// Result is what the channel returns
type Result struct {
	Msg   uint64
	Error error
}

// Init initializes the consensus protocol
func (r *RBC) Init(c *config.NodeConfig) {
	r.isSender = c.Config.ProtConfig.Id == SenderID
}

// Setup sets up the network component
func (r *RBC) Setup(n *net.Network) {
	r.host = n.H
	r.ctx = n.Ctx
	r.pMap = n.PeerMap
	r.host.SetStreamHandler(ID, r.MsgHandler)
}

// MsgHandler talks to/from the network
func (r *RBC) MsgHandler(s network.Stream) {
	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		select {
		case err, ok := <-errCh:
			if ok {
				fmt.Println(err)
			} else {
				fmt.Println("Some error")
			}
		}
		s.Reset()
	}()

	buf := make([]byte, 1)
	// Event Handler
	for {
		// Receive a value from the sender and print it
		_, err := io.ReadFull(s, buf)
		if err != nil {
			errCh <- err
			return
		}
		fmt.Println("Got from sender a value", buf)
	}
}

// Broadcast broadcasts a byte to all the nodes
func (r *RBC) Broadcast(b byte) error {
	for idx, p := range r.pMap {
		fmt.Println("Talking to node", idx, "Peer: ", p)
		s, err := r.host.NewStream(r.ctx, p.ID, ID)
		if err != nil {
			return err
		}
		_, err = s.Write([]byte{b})
		if err != nil {
			return err
		}
	}
	return nil
}

// Start implements the Protocol interface
func (r *RBC) Start() {
	if r.isSender {
		err := r.Broadcast(byte(10))
		if err != nil {
			panic(err)
		}
	}
	// Wait a minute
	m1, _ := time.ParseDuration("1m")
	<-time.After(m1)
}
