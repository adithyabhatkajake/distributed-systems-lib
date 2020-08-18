package main

/*
 * A client does the following:
 * Read the config to get public key and IP maps
 * Let B be the number of commands.
 * Send B commands to the nodes and wait for f+1 acknowledgements
 */

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

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

type vote struct {
	voter uint64
	cmd   crypto.Hash
}

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
	voteChannel     chan vote
	idMap           = make(map[string]uint64)
	votes           = make(map[crypto.Hash]uint64)
	f               uint64
)

func sendRequest(cmdChannel chan *msg.E2CMsg, rwMap map[uint64]*bufio.Writer) {
	for {
		cmdMsg, ok := <-cmdChannel
		fmt.Println("Processing Command")
		if !ok {
			fmt.Println("sendRequest channel processing error")
			fmt.Println("Closing thread")
			return
		}
		data, err := pb.Marshal(cmdMsg)
		if err != nil {
			fmt.Println("Marshaling error", err)
			continue
		}
		// Ship command off to the nodes
		for idx, rw := range rwMap {
			streamMutex.Lock()
			rw.Write(data)
			rw.Flush()
			streamMutex.Unlock()
			fmt.Println("Sending command to node", idx)
		}
	}
}

func ackMsgHandler(s network.Stream) {
	reader := bufio.NewReader(s)
	// Prepare a buffer to receive an acknowledgement
	msgBuf := make([]byte, msg.MaxMsgSize)
	for {
		_, err := reader.Read(msgBuf)
		if err != nil {
			fmt.Println("bufio read error", err)
			return
		}
		// Ensure Correct Signature
		v := vote{
			voter: idMap[s.ID()],
		}
		voteChannel <- v
	}
}

func handleVotes() {
	for {
		// Get Acknowledgements from nodes after consensus
		v, ok := <-voteChannel
		if !ok {
			fmt.Println("vote channel closed")
			return
		}
		// Deal with the vote
		// TODO
		fmt.Println("Received a vote from node", v.voter)
		voteMutex.Lock()
		votes[v.cmd]++
		if votes[v.cmd] > f {
			fmt.Println(v.cmd, "is committed.")
		}
		voteMutex.Unlock()
		// We got acknowledgement for a command
		cmdMutex.Lock()
		PendingCommands--
		cmdMutex.Unlock()
	}
}

func main() {
	fmt.Println("I am the client")
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
	fmt.Println("Client at", node.Addrs())

	pMap := make(map[uint64]peerstore.AddrInfo)
	streamMap := make(map[uint64]network.Stream)
	rwMap := make(map[uint64]*bufio.Writer)
	connectedNodes := uint64(0)

	for i := uint64(0); i < confData.GetNumNodes(); i++ {
		// Prepare peerInfo
		pMap[i] = confData.GetPeerFromID(i)
		// Connect to node i
		fmt.Println("Attempting connection to node", pMap[i])
		err = node.Connect(ctx, pMap[i])
		if err != nil {
			fmt.Println("Connection Error", err)
			continue
		}
		streamMap[i], err = node.NewStream(ctx, pMap[i].ID,
			e2cconsensus.ClientProtocolID)
		if err != nil {
			fmt.Println("Stream opening Error", err)
			continue
		}
		idMap[streamMap[i].ID()] = i
		connectedNodes++
		rwMap[i] = bufio.NewWriter(streamMap[i])
	}
	// Handle all messages received here
	node.SetStreamHandler(e2cconsensus.ClientProtocolID, ackMsgHandler)

	// Ensure we are connected to sufficient nodes
	if connectedNodes <= f {
		fmt.Println("Insufficient connections to replicas")
		return
	}

	cmdChannel := make(chan *msg.E2CMsg, BufferCommands)
	voteChannel = make(chan vote, BufferCommands)

	// Run a goroutine that keeps sending requests to the nodes
	go sendRequest(cmdChannel, rwMap)

	// Spawn a thread that handles acknowledgement received for the
	// various requests
	go handleVotes()

	// Make sure we always fill the channel with commands
	idx := uint64(0)
	for {
		skipLoop := false
		cmdMutex.Lock()
		if PendingCommands >= BufferCommands {
			skipLoop = true
		} else {
			PendingCommands++
		}
		cmdMutex.Unlock()

		// If the buffer is full, skip sending a command
		if skipLoop {
			continue
		}

		// Increment index
		idx++

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

		fmt.Println("Sending command to thread")
		// Dispatch E2C message for processing
		cmdChannel <- cmdMsg
	}
}
