package net

import (
	"context"
	"sync"
	"time"

	"github.com/adithyabhatkajake/libchatter/log"

	"github.com/adithyabhatkajake/libchatter/config"
	"github.com/adithyabhatkajake/libchatter/crypto"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

var (
	// RetryLimit specifies how many times to try dialing a node
	RetryLimit = 30
	// RetryWaitDuration specifies how many times to wait between each tries
	RetryWaitDuration = "1s"
)

// Network contains all the networking related data structures
type Network struct {
	H          host.Host
	Ctx        context.Context
	CancelFunc context.CancelFunc
	PeerMap    map[uint64]*peerstore.AddrInfo
}

/* 	Write a package that implements the following.
1. Setup : takes a config and connects to all the clients
2. Send: sends a message to a specific recipient
3. Broadcast: Sends a message to all recipients
*/

// Setup function sets up connection with all the nodes
// Setup needs to start listening on its port and try connecting to other nodes
func Setup(c config.Config, nc Conf, cryptoc crypto.Config) *Network {
	var err error
	ctx := context.Background()
	net := &Network{}
	net.Ctx, net.CancelFunc = context.WithCancel(ctx)
	myID := nc.GetID()
	net.H, err = libp2p.New(net.Ctx,
		libp2p.ListenAddrStrings(nc.GetP2PAddrFromID(myID)),
		libp2p.Identity(cryptoc.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}
	numNodes := c.GetNumNodes()
	peerMap := make(map[uint64]*peerstore.AddrInfo)
	for idx := uint64(0); idx < numNodes; idx++ {
		if idx == myID {
			continue
		}
		peerID, err := peerstore.IDFromPublicKey(cryptoc.GetPubKeyFromID(idx))
		if err != nil {
			panic(err)
		}
		addr, err := ma.NewMultiaddr(nc.GetP2PAddrFromID(idx))
		if err != nil {
			panic(err)
		}
		peerMap[idx] = &peerstore.AddrInfo{
			ID:    peerID,
			Addrs: []ma.Multiaddr{addr},
		}
	}
	net.PeerMap = peerMap
	return net
}

// Connect connects to all the peers in parallel.
// Ensure that Setup is called before calling this
func (n *Network) Connect() {
	wg := &sync.WaitGroup{}
	for _, peer := range n.PeerMap {
		wg.Add(1)
		go n.connectPeer(wg, peer)
	}
	wg.Wait()
}

func (n *Network) connectPeer(wg *sync.WaitGroup, p *peerstore.AddrInfo) {
	defer wg.Done()
	var err error
	t, _ := time.ParseDuration(RetryWaitDuration)
	for i := 0; i < RetryLimit; i++ {
		log.Info("Attempting connection to", *p)
		err = n.H.Connect(n.Ctx, *p)
		if err != nil {
			log.Debug("Connection Failed. Retrying Attempt #", i)
			<-time.After(t)
			continue
		}
		break
	}
	// If we still fail to connect, panic
	if err != nil {
		panic(err)
	}
}

// ShutDown closes the node and disconnects from the network
func (n *Network) ShutDown() {
	defer n.CancelFunc()
	err := n.H.Close()
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
