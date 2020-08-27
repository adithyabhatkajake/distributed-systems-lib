package synchs

import (
	"bufio"
	"context"
	"sync"
	"time"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/log"
	"github.com/adithyabhatkajake/libe2c/util"

	"github.com/libp2p/go-libp2p"

	config "github.com/adithyabhatkajake/libe2c/config/synchs"
	"github.com/adithyabhatkajake/libe2c/net"

	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"

	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

const (
	// ProtocolID is the ID for E2C Protocol
	ProtocolID = "synchs/synchs/0.0.1"
	// ProtocolMsgBuffer defines how many protocol messages can be buffered
	ProtocolMsgBuffer = 100
)

// Init implements the Protocol interface
func (shs *SyncHS) Init(c *config.NodeConfig) {
	shs.config = c
	shs.leader = DefaultLeaderID
	shs.view = 1 // View Number starts from 1
	shs.blTimer = nil
}

// Setup sets up the network components
func (shs *SyncHS) Setup(n *net.Network) error {
	shs.host = n.H
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(shs.config.GetClientListenAddr()),
		libp2p.Identity(shs.config.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}
	shs.cliHost = host
	shs.ctx = n.Ctx

	// Setup maps
	shs.pMap = n.PeerMap
	shs.streamMap = make(map[uint64]*bufio.ReadWriter)
	shs.cliMap = make(map[*bufio.ReadWriter]bool)
	shs.pendingCommands = make(map[crypto.Hash]*chain.Command)
	shs.timerMaps = make(map[uint64]*util.Timer)
	shs.blameMap = make(map[uint64]map[uint64]*msg.Blame)
	shs.certMap = make(map[uint64]*msg.BlockCertificate)

	// Setup channels
	shs.msgChannel = make(chan *msg.SyncHSMsg, ProtocolMsgBuffer)
	shs.cmdChannel = make(chan *chain.Command, ProtocolMsgBuffer)
	shs.voteChannel = make(chan *msg.Vote, ProtocolMsgBuffer)

	// Obtain a new chain
	shs.bc = chain.NewChain()
	// TODO: create a new chain only if no chain is present in the data directory

	// Setup certificate for the first block
	genesisCert := &msg.BlockCertificate{
		BCert: &msg.Certificate{},
		Data:  &msg.VoteData{},
	}
	genesisCert.Data.View = shs.view
	genesisCert.Data.Block = chain.GetGenesis()
	shs.certMap[0] = genesisCert

	// How to react to Protocol Messages
	shs.host.SetStreamHandler(ProtocolID, shs.ProtoMsgHandler)

	// How to react to Client Messages
	shs.cliHost.SetStreamHandler(ClientProtocolID, shs.ClientMsgHandler)

	// Connect to all the other nodes talking E2C protocol
	wg := &sync.WaitGroup{} // For faster setup
	for idx, p := range shs.pMap {
		wg.Add(1)
		go func(idx uint64, p *peerstore.AddrInfo) {
			log.Trace("Attempting to open a stream with", p, "using protocol", ProtocolID)
			retries := 300
			for i := retries; i > 0; i-- {
				s, err := shs.host.NewStream(shs.ctx, p.ID, ProtocolID)
				if err != nil {
					log.Error("Error connecting to peers:", err)
					log.Info("Retry attempt ", retries-i+1, " to connect to node ", idx, " in a second")
					<-time.After(time.Second)
					continue
				}
				shs.netMutex.Lock()
				shs.streamMap[idx] = bufio.NewReadWriter(
					bufio.NewReader(s), bufio.NewWriter(s))
				shs.netMutex.Unlock()
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
func (shs *SyncHS) Start() {
	// First, start vote handler concurrently
	go shs.voteHandler()
	// Then start command handler concurrently
	go shs.cmdHandler()
	// Start E2C Protocol - Start message handler
	shs.protocol()
}

// ProtoMsgHandler reacts to all protocol messages in the network
func (shs *SyncHS) ProtoMsgHandler(s network.Stream) {
	shs.errCh = make(chan error, 1)
	defer close(shs.errCh)

	// Run a parallel goroutine that closes the channel if any error is detected
	go func() {
		select {
		case err, ok := <-shs.errCh:
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
			shs.errCh <- err
			return
		}
		// Use a copy of the message and send it to off for processing
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React to the message in parallel and continue
		go shs.react(msgBuf)
	}
}
