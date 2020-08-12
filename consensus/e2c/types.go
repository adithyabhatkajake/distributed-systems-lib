package e2c

import (
	"context"

	"github.com/libp2p/go-libp2p-core/network"

	chain "github.com/adithyabhatkajake/libe2c/chain"
	config "github.com/adithyabhatkajake/libe2c/config/e2c"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"

	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

// E2C implements the consensus protocol
type E2C struct {
	// network data structures
	host       host.Host
	cliHost    host.Host
	ctx        context.Context
	pMap       map[uint64]*peerstore.AddrInfo
	streamMap  map[uint64]network.Stream
	msgChannel chan msg.E2CMsg
	bc         chain.BlockChain

	// Protocol information
	leader uint64
	config *config.NodeConfig

	// Error Channel
	errCh chan error
}
