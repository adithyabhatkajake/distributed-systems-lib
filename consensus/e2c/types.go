package e2c

import (
	"bufio"
	"context"
	"sync"

	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/util"

	chain "github.com/adithyabhatkajake/libe2c/chain"
	config "github.com/adithyabhatkajake/libe2c/config/e2c"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"

	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

// E2C implements the consensus protocol
type E2C struct {
	// Network data structures
	host    host.Host
	cliHost host.Host
	ctx     context.Context

	// Maps
	pMap            map[uint64]*peerstore.AddrInfo
	streamMap       map[uint64]*bufio.ReadWriter
	pendingCommands map[crypto.Hash]*chain.Command
	timerMaps       map[uint64]*util.Timer

	// Locks
	cmdMutex sync.Mutex
	netMutex sync.Mutex

	// Channels
	msgChannel chan *msg.E2CMsg
	cmdChannel chan *chain.Command
	errCh      chan error

	// Block chain
	bc *chain.BlockChain

	// Protocol information
	leader uint64
	config *config.NodeConfig
}
