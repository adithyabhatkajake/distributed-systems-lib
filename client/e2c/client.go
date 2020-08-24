package main

/*
 * A client does the following:
 * Read the config to get public key and IP maps
 * Let B be the number of commands.
 * Send B commands to the nodes and wait for f+1 acknowledgements for every acknowledgement
 */

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/adithyabhatkajake/libe2c/log"
	"github.com/adithyabhatkajake/libe2c/util"

	"github.com/adithyabhatkajake/libe2c/chain"
	"github.com/adithyabhatkajake/libe2c/config/e2c"
	e2cconsensus "github.com/adithyabhatkajake/libe2c/consensus/e2c"
	"github.com/adithyabhatkajake/libe2c/crypto"
	e2cio "github.com/adithyabhatkajake/libe2c/io"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"

	pb "github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

var (
	// BufferCommands defines how many commands to wait for
	// acknowledgement in a batch
	BufferCommands = uint64(10)
	// PendingCommands tells how many commands we are waiting
	// for acknowledgements from replicas
	PendingCommands = uint64(0)
	cmdMutex        = &sync.Mutex{}
	streamMutex     = &sync.Mutex{}
	voteMutex       = &sync.Mutex{}
	condLock        = &sync.Mutex{}
	cond            = sync.NewCond(condLock)
	voteChannel     chan *msg.CommitAck
	idMap           = make(map[string]uint64)
	votes           = make(map[crypto.Hash]uint64)
	f               uint64
)

func sendCommandToServer(cmd *msg.E2CMsg, rwMap map[uint64]*bufio.Writer) {
	log.Trace("Processing Command")
	data, err := pb.Marshal(cmd)
	if err != nil {
		log.Error("Marshaling error", err)
		return
	}
	// Ship command off to the nodes
	for idx, rw := range rwMap {
		streamMutex.Lock()
		rw.Write(data)
		rw.Flush()
		streamMutex.Unlock()
		log.Trace("Sending command to node", idx)
	}
}

func ackMsgHandler(s network.Stream, serverID uint64) {
	reader := bufio.NewReader(s)
	// Prepare a buffer to receive an acknowledgement
	msgBuf := make([]byte, msg.MaxMsgSize)
	log.Trace("Started acknowledgement message handler")
	for {
		len, err := reader.Read(msgBuf)
		if err != nil {
			log.Error("bufio read error", err)
			return
		}
		log.Trace("Received a message from the server.", serverID)
		msg := &msg.E2CMsg{}
		err = pb.Unmarshal(msgBuf[0:len], msg)
		if err != nil {
			log.Error("Unmarshalling error", serverID, err)
			continue
		}
		voteChannel <- msg.GetAck()
	}
}

func handleVotes(cmdChannel chan *msg.E2CMsg, rwMap map[uint64]*bufio.Writer) {
	voteMap := make(map[crypto.Hash]uint64)
	commitMap := make(map[crypto.Hash]bool)
	for {
		// Get Acknowledgements from nodes after consensus
		v, ok := <-voteChannel
		log.Trace("Received an acknowledgement")
		if !ok {
			log.Error("vote channel closed")
			return
		}
		if v == nil {
			continue
		}
		// Deal with the vote
		// We received a vote. Now add this to the conformation map
		cmdHash := crypto.ToHash(v.CmdHash)
		_, exists := voteMap[cmdHash]
		if !exists {
			voteMap[cmdHash] = 1       // 1 means we have seen one vote so far.
			commitMap[cmdHash] = false // To say that we have not yet committed this value
		} else {
			voteMap[cmdHash]++ // Add another vote
		}
		// To ensure this is executed only once, check old committed state
		old := commitMap[cmdHash]
		if voteMap[cmdHash] > f {
			log.Info("Confirmed commit for command ",
				util.HashToString(cmdHash))
			commitMap[cmdHash] = true
		}
		new := commitMap[cmdHash]
		// If we commit the block for the first time, then ship off a new command to the server
		if old != new {
			cmd := <-cmdChannel
			log.Info("Sending command ", cmd, " to the servers")
			go sendCommandToServer(cmd, rwMap)
		}
	}
}

func main() {
	// Setup Logger
	log.SetLevel(log.InfoLevel)

	log.Info("I am the client")
	ctx := context.Background()

	// Get client config
	confData := &e2c.ClientConfig{}
	e2cio.ReadFromFile(confData, os.Args[1])

	f = uint64((confData.GetNumNodes() - 1) / 2)
	// Start networking stack
	node, err := p2p.New(ctx,
		libp2p.Identity(confData.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}

	// Print self information
	log.Info("Client at", node.Addrs())

	// Handle all messages received using ackMsgHandler
	// node.SetStreamHandler(e2cconsensus.ClientProtocolID, ackMsgHandler)
	// Setting stream handler is useless :/

	pMap := make(map[uint64]peerstore.AddrInfo)
	streamMap := make(map[uint64]network.Stream)
	rwMap := make(map[uint64]*bufio.Writer)
	connectedNodes := uint64(0)

	for i := uint64(0); i < confData.GetNumNodes(); i++ {
		// Prepare peerInfo
		pMap[i] = confData.GetPeerFromID(i)
		// Connect to node i
		log.Trace("Attempting connection to node ", pMap[i])
		err = node.Connect(ctx, pMap[i])
		if err != nil {
			log.Error("Connection Error ", err)
			continue
		}
		streamMap[i], err = node.NewStream(ctx, pMap[i].ID,
			e2cconsensus.ClientProtocolID)
		if err != nil {
			log.Error("Stream opening Error", err)
			continue
		}
		idMap[streamMap[i].ID()] = i
		connectedNodes++
		rwMap[i] = bufio.NewWriter(streamMap[i])
		go ackMsgHandler(streamMap[i], i)
	}

	// Ensure we are connected to sufficient nodes
	if connectedNodes <= f {
		log.Warn("Insufficient connections to replicas")
		return
	}

	cmdChannel := make(chan *msg.E2CMsg, BufferCommands)
	voteChannel = make(chan *msg.CommitAck, BufferCommands)

	// First, spawn a thread that handles acknowledgement received for the
	// various requests
	go handleVotes(cmdChannel, rwMap)

	idx := uint64(0)

	// Then, run a goroutine that sends the first BufferCommands requests to the nodes
	for ; idx < BufferCommands; idx++ {
		// Send via stream a command
		cmdStr := fmt.Sprintf("Do my bidding #%d my servant!", idx)

		// Build a command
		cmd := &chain.Command{}
		// Set command
		cmd.Cmd = []byte(cmdStr)
		// Sign the command
		cmd.Clientsig, err = confData.GetMyKey().Sign(cmd.Cmd)
		if err != nil {
			panic(err)
		}

		// Build a protocol message
		cmdMsg := &msg.E2CMsg{}
		cmdMsg.Msg = &msg.E2CMsg_Cmd{
			Cmd: cmd,
		}

		log.Info("Sending command ", idx, " to the servers")
		go sendCommandToServer(cmdMsg, rwMap)
	}

	// Make sure we always fill the channel with commands

	for {
		// Send via stream a command
		cmdStr := fmt.Sprintf("Do my bidding #%d my servant!", idx)

		// Build a command
		cmd := &chain.Command{}
		// Set command
		cmd.Cmd = []byte(cmdStr)
		// Sign the command
		cmd.Clientsig, err = confData.GetMyKey().Sign(cmd.Cmd)
		if err != nil {
			panic(err)
		}

		// Build a protocol message
		cmdMsg := &msg.E2CMsg{}
		cmdMsg.Msg = &msg.E2CMsg_Cmd{
			Cmd: cmd,
		}

		// Dispatch E2C message for processing
		// This will block until some of the commands are committed
		cmdChannel <- cmdMsg
		// Increment command number
		idx++
	}
}
