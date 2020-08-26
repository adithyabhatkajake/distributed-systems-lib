package e2c

import (
	"fmt"
	"time"

	"github.com/adithyabhatkajake/libe2c/crypto"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// Implement all the interfaces, i.e.,
// 1. net
// 2. crypto
// 3. config

// GetID returns the Id of this instance
func (ec *NodeConfig) GetID() uint64 {
	return ec.Config.ProtConfig.Id
}

// GetP2PAddrFromID gets the P2P address of the node rid
func (ec *NodeConfig) GetP2PAddrFromID(rid uint64) string {
	address := ec.Config.NetConfig.NodeAddressMap[rid]
	addr := fmt.Sprintf("/ip4/%s/tcp/%s", address.IP, address.Port)
	return addr
}

// GetP2PAddrFromID gets the P2P address of the node rid
func (cc *ClientConfig) GetP2PAddrFromID(rid uint64) string {
	address := cc.Config.NetConfig.NodeAddressMap[rid]
	addr := fmt.Sprintf("/ip4/%s/tcp/%s", address.IP, address.Port)
	return addr
}

// GetMyKey returns the private key of this instance
func (ec *NodeConfig) GetMyKey() crypto.PrivKey {
	return ec.PvtKey
}

// GetMyKey returns the private key of this instance
func (cc *ClientConfig) GetMyKey() crypto.PrivKey {
	return cc.PvtKey
}

// GetPubKeyFromID returns the Public key of node whose ID is nid
func (ec *NodeConfig) GetPubKeyFromID(nid uint64) crypto.PubKey {
	return ec.NodeKeyMap[nid]
}

// GetPubKeyFromID returns the Public key of node whose ID is nid
func (cc *ClientConfig) GetPubKeyFromID(nid uint64) crypto.PubKey {
	return cc.NodeKeyMap[nid]
}

// GetPeerFromID returns libp2p peerInfo from the config
func (ec *NodeConfig) GetPeerFromID(nid uint64) peerstore.AddrInfo {
	pID, err := peerstore.IDFromPublicKey(ec.GetPubKeyFromID(0))
	if err != nil {
		panic(err)
	}
	addr, err := ma.NewMultiaddr(ec.GetP2PAddrFromID(0))
	if err != nil {
		panic(err)
	}
	pInfo := peerstore.AddrInfo{
		ID:    pID,
		Addrs: []ma.Multiaddr{addr},
	}
	return pInfo
}

// GetPeerFromID returns libp2p peerInfo from the config
func (cc *ClientConfig) GetPeerFromID(nid uint64) peerstore.AddrInfo {
	pID, err := peerstore.IDFromPublicKey(cc.GetPubKeyFromID(nid))
	if err != nil {
		panic(err)
	}
	addr, err := ma.NewMultiaddr(cc.GetP2PAddrFromID(nid))
	if err != nil {
		panic(err)
	}
	pInfo := peerstore.AddrInfo{
		ID:    pID,
		Addrs: []ma.Multiaddr{addr},
	}
	return pInfo
}

// GetNumNodes returns the protocol size
func (ec *NodeConfig) GetNumNodes() uint64 {
	return ec.Config.ProtConfig.Info.NodeSize
}

// GetNumNodes returns the protocol size
func (cc *ClientConfig) GetNumNodes() uint64 {
	return cc.Config.Info.NodeSize
}

// GetClientListenAddr returns the address where to talk to/from clients
func (ec *NodeConfig) GetClientListenAddr() string {
	id := ec.GetID()
	address := ec.Config.ClientNetConfig.NodeAddressMap[id]
	addr := fmt.Sprintf("/ip4/%s/tcp/%s", address.IP, address.Port)
	return addr
}

// GetBlockSize returns the number of commands that can be inserted in one block
func (ec *NodeConfig) GetBlockSize() uint64 {
	return ec.Config.GetProtConfig().GetInfo().GetBlockSize()
}

// GetDelta returns the synchronous wait time
func (ec *NodeConfig) GetDelta() time.Duration {
	timeInSeconds := ec.Config.ProtConfig.GetDelta()
	return time.Duration(int(timeInSeconds*1000)) * time.Millisecond
}

// GetCommitWaitTime returns how long to wait before committing a block
func (ec *NodeConfig) GetCommitWaitTime() time.Duration {
	return ec.GetDelta() * 2
}

// GetNPBlameWaitTime returns how long to wait before sending the NP Blame
func (ec *NodeConfig) GetNPBlameWaitTime() time.Duration {
	return ec.GetDelta() * 3
}

// GetNumberOfFaultyNodes computes f for this protocol as f = (n-1)/2
func (ec *NodeConfig) GetNumberOfFaultyNodes() uint64 {
	return uint64((ec.GetNumNodes() - 1) / 2)
}
