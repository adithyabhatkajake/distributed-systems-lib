package e2c

import (
	"context"

	"github.com/libp2p/go-libp2p-core/network"

	"github.com/adithyabhatkajake/libe2c/config"

	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

// E2C implements the consensus protocol
type E2C struct {
	// network data structures
	host      host.Host
	ctx       context.Context
	pMap      map[uint64]*peerstore.AddrInfo
	streamMap map[uint64]network.Stream

	// Protocol information
	leader     uint64
	protConfig *config.E2CConfig

	// Error Channel
	errCh chan error
}
