package main

import (
	"flag"
	"os"

	config "github.com/adithyabhatkajake/libe2c/config/apollo"
	synchs "github.com/adithyabhatkajake/libe2c/consensus/apollo"
	"github.com/adithyabhatkajake/libe2c/io"
	"github.com/adithyabhatkajake/libe2c/log"
	"github.com/adithyabhatkajake/libe2c/net"
)

func main() {
	// Parse flags
	logLevelPtr := flag.Uint64("loglevel", uint64(log.InfoLevel),
		"Loglevels are one of \n0 - PanicLevel\n1 - FatalLevel\n2 - ErrorLevel\n3 - WarnLevel\n4 - InfoLevel\n5 - DebugLevel\n6 - TraceLevel")
	configFileStrPtr := flag.String("conf", "", "Path to config file")

	flag.Parse()

	logLevel := log.InfoLevel

	switch uint32(*logLevelPtr) {
	case 0:
		logLevel = log.PanicLevel
	case 1:
		logLevel = log.FatalLevel
	case 2:
		logLevel = log.ErrorLevel
	case 3:
		logLevel = log.WarnLevel
	case 4:
		logLevel = log.InfoLevel
	case 5:
		logLevel = log.DebugLevel
	case 6:
		logLevel = log.TraceLevel
	}

	// Log Settings
	log.SetLevel(logLevel)

	log.Info("I am the replica.")
	Config := &config.NodeConfig{}

	io.ReadFromFile(Config, *configFileStrPtr)
	log.Debug("Finished reading the config file", os.Args[1])

	// Setup connections
	netw := net.Setup(Config, Config, Config)

	// Connect and send a test message
	netw.Connect()
	log.Debug("Finished connection to all the nodes")

	// Configure E2C protocol
	n := &synchs.Apollo{}
	n.Init(Config)
	n.Setup(netw)

	// Start E2C
	n.Start()

	// Disconnect
	netw.ShutDown()
}
