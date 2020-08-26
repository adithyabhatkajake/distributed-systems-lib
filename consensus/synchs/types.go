package synchs

import (
	"bufio"
	"context"
	"sync"

	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/util"

	chain "github.com/adithyabhatkajake/libe2c/chain"
	config "github.com/adithyabhatkajake/libe2c/config/synchs"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"

	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

// SyncHS implements the consensus protocol
type SyncHS struct {
	// Network data structures
	host    host.Host
	cliHost host.Host
	ctx     context.Context

	// Maps
	// Mapping between ID and libp2p-peer
	pMap map[uint64]*peerstore.AddrInfo
	// A set of all known clients
	cliMap map[*bufio.ReadWriter]bool
	// A map of node ID to its corresponding RW stream
	streamMap map[uint64]*bufio.ReadWriter
	// A map of hash to pending commands
	pendingCommands map[crypto.Hash]*chain.Command
	// A mapping between the block number to its commit timer
	timerMaps map[uint64]*util.Timer
	// A mapping between the view and (A mapping between the origin and blames against the leader)
	blameMap map[uint64]map[uint64]*msg.Blame

	/* Locks - We separate all the locks, so that acquiring
	one lock does not make other goroutines stop */
	peerMapLock sync.RWMutex // The lock to modify
	cliMutex    sync.RWMutex // The lock to modify cliMap
	netMutex    sync.RWMutex // The lock to modify streamMap
	cmdMutex    sync.RWMutex // The lock to modify pendingCommands
	timerLock   sync.RWMutex // The lock to modify timerMaps
	blTimerLock sync.RWMutex // The lock to modify blTimer
	blLock      sync.RWMutex // The lock to modify blameMap

	// Channels
	msgChannel chan *msg.SyncHSMsg
	cmdChannel chan *chain.Command
	errCh      chan error

	// Block chain
	bc *chain.BlockChain

	// Protocol information
	leader  uint64
	view    uint64
	config  *config.NodeConfig
	blTimer *util.Timer
}
