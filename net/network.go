package net

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/adithyabhatkajake/libe2c/config"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

var (
	// RetryLimit specifies how many times to try dialing a node
	RetryLimit = 30
	// RetryWaitDuration specifies how many times to wait between each tries
	RetryWaitDuration = "20s"
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
func Setup(nodeConf *config.NodeConfig) *Network {
	var err error
	ctx := context.Background()
	net := &Network{}
	net.Ctx, net.CancelFunc = context.WithCancel(ctx)
	myID := nodeConf.Config.ProtConfig.Id
	net.H, err = libp2p.New(net.Ctx,
		libp2p.ListenAddrStrings(
			nodeConf.Config.NetConfig.NodeAddressMap[myID].GetP2PAddr()),
		libp2p.Identity(nodeConf.PvtKey),
	)
	if err != nil {
		panic(err)
	}
	peerMap := make(map[uint64]*peerstore.AddrInfo)
	for idx, Addr := range nodeConf.Config.NetConfig.NodeAddressMap {
		if idx == myID {
			continue
		}
		peerID, err := peerstore.IDFromPublicKey(nodeConf.NodeKeyMap[idx])
		if err != nil {
			panic(err)
		}
		addr, err := ma.NewMultiaddr(Addr.GetP2PAddr())
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

// Connect connects to all the peers.
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
		err = n.H.Connect(n.Ctx, *p)
		if err != nil {
			fmt.Println("Connection Failed. Retrying Attempt #", i)
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
		panic(err)
	}
}
