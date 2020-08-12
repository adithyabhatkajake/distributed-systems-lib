package main

import (
	"fmt"
	"os"

	config "github.com/adithyabhatkajake/libe2c/config/e2c"
	"github.com/adithyabhatkajake/libe2c/consensus/e2c"
	"github.com/adithyabhatkajake/libe2c/io"
	"github.com/adithyabhatkajake/libe2c/net"
)

func main() {
	fmt.Println("I am the replica.")
	Config := &config.NodeConfig{}

	io.ReadFromFile(Config, os.Args[1])
	fmt.Println("Finished reading the config file", os.Args[1])

	// Setup connections
	netw := net.Setup(Config, Config, Config)

	// Connect and send a test message
	netw.Connect()
	fmt.Println("Finished connection to all the nodes")

	// Configure E2C protocol
	e := &e2c.E2C{}
	e.Init(Config)
	e.Setup(netw)

	// Start E2C
	e.Start()

	// Disconnect
	netw.ShutDown()
}
