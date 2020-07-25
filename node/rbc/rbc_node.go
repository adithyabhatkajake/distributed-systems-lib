package main

/* We implement the replica portion of the protocol here. */

import (
	"fmt"
	"os"

	"github.com/adithyabhatkajake/libe2c/config"
	"github.com/adithyabhatkajake/libe2c/consensus"
	"github.com/adithyabhatkajake/libe2c/io"
	"github.com/adithyabhatkajake/libe2c/net"
)

var (
	waitTime = "20s"
)

func main() {
	fmt.Println("I am the replica.")
	config := &config.NodeConfig{}

	io.ReadFromFile(config, os.Args[1])
	fmt.Println("Finished reading the config file", os.Args[1])

	// Setup connections
	netw := net.Setup(config)

	// Connect and send a test message
	netw.Connect()
	fmt.Println("Finished connection to all the nodes")

	// Configure RBC protocol
	r := &consensus.RBC{}
	r.Init(config)
	r.Setup(netw)

	// Start RBC
	r.Start()

	// Disconnect
	netw.ShutDown()
}
