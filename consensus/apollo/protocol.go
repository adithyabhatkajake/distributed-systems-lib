package apollo

import (
	"bufio"
	"context"
	"sync"
	"time"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/log"

	"github.com/libp2p/go-libp2p"

	config "github.com/adithyabhatkajake/libe2c/config/apollo"
	"github.com/adithyabhatkajake/libe2c/net"

	msg "github.com/adithyabhatkajake/libe2c/msg/apollo"

	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

const (
	// ProtocolID is the ID for E2C Protocol
	ProtocolID = "apollo/apollo/0.0.1"
	// ProtocolMsgBuffer defines how many protocol messages can be buffered
	ProtocolMsgBuffer = 100
)

// Init implements the Protocol interface
func (n *Apollo) Init(c *config.NodeConfig) {
	n.config = c
	n.leader = DefaultLeaderID
	n.blTimer = nil
}

// Setup sets up the network components
func (n *Apollo) Setup(netw *net.Network) error {
	n.host = netw.H
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(n.config.GetClientListenAddr()),
		libp2p.Identity(n.config.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}
	n.cliHost = host
	n.ctx = netw.Ctx

	// Setup maps
	n.pMap = netw.PeerMap
	n.streamMap = make(map[uint64]*bufio.ReadWriter)
	n.cliMap = make(map[*bufio.ReadWriter]bool)
	n.pendingCommands = make(map[crypto.Hash]*chain.Command)
	n.blameMap = make(map[uint64]map[uint64]*msg.Blame)

	// Setup channels
	n.msgChannel = make(chan *msg.ApolloMsg, ProtocolMsgBuffer)
	n.cmdChannel = make(chan *chain.Command, ProtocolMsgBuffer)

	// Obtain a new chain
	n.bc = chain.NewChain()
	// TODO: create a new chain only if no chain is present in the data directory

	// How to react to Protocol Messages
	n.host.SetStreamHandler(ProtocolID, n.ProtoMsgHandler)

	// How to react to Client Messages
	n.cliHost.SetStreamHandler(ClientProtocolID, n.ClientMsgHandler)

	// Connect to all the other nodes talking E2C protocol
	wg := &sync.WaitGroup{} // For faster setup
	for idx, p := range n.pMap {
		wg.Add(1)
		go func(idx uint64, p *peerstore.AddrInfo) {
			log.Trace("Attempting to open a stream with", p, "using protocol", ProtocolID)
			retries := 300
			for i := retries; i > 0; i-- {
				s, err := n.host.NewStream(n.ctx, p.ID, ProtocolID)
				if err != nil {
					log.Error("Error connecting to peers:", err)
					log.Info("Retry attempt ", retries-i+1, " to connect to node ", idx, " in a second")
					<-time.After(time.Second)
					continue
				}
				n.netMutex.Lock()
				n.streamMap[idx] = bufio.NewReadWriter(
					bufio.NewReader(s), bufio.NewWriter(s))
				n.netMutex.Unlock()
				log.Info("Connected to Node ", idx)
				break
			}
			wg.Done()
		}(idx, p)
	}
	wg.Wait()
	log.Info("Setup Finished. Ready to do SMR:)")

	return nil
}

// Start implements the Protocol Interface
func (n *Apollo) Start() {
	// Concurrently handle commands
	go n.cmdHandler()
	// Start E2C Protocol - Start message handler
	n.protocol()
}

// ProtoMsgHandler reacts to all protocol messages in the network
func (n *Apollo) ProtoMsgHandler(s network.Stream) {
	n.errCh = make(chan error, 1)
	defer close(n.errCh)

	// Run a parallel goroutine that closes the channel if any error is detected
	go func() {
		select {
		case err, ok := <-n.errCh:
			if ok {
				log.Error("Type-I ", err, ok)
			} else {
				log.Error("Type-II", err, ok)
			}
		}
		s.Close()
	}()

	// A global buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	// Event Handler
	reader := bufio.NewReader(s)
	for {
		// Receive a message from anyone and process them
		len, err := reader.Read(buf)
		if err != nil {
			n.errCh <- err
			return
		}
		// Use a copy of the message and send it to off for processing
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React to the message in parallel and continue
		go n.react(msgBuf)
	}
}
