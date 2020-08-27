package main

import (
	"fmt"
	"os"

	"github.com/adithyabhatkajake/libe2c/io"

	"github.com/pborman/getopt"

	"github.com/adithyabhatkajake/libe2c/config"
	e2cconfig "github.com/adithyabhatkajake/libe2c/config/e2c"
	"github.com/adithyabhatkajake/libe2c/crypto/secp256k1"
)

var (
	nReplicas      uint64  = 10
	nFaulty        uint64  = 4
	blkSize        uint64  = 1
	delta          float64 = 1 // in seconds
	basePort       uint32  = 10000
	clientBasePort uint32  = 20000
	outDir         string  = "testData/"
	defaultIP      string  = "127.0.0.1"
)

var (
	generalConfigFile = "nodes.txt"
	nodeConfigFile    = "nodes-%d.txt"
	clientFile        = "client.txt"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/* This program performs the followins things: */
/* 1. Generate n public-private key pairs */
/* 2. Print the <replicaID, address, port, public key> in nodes.txt */
/* 3. Print the <Private Key> in nodes-<nodeID>.txt for every node */
func main() {
	optNumReplica := getopt.Uint64('n', nReplicas, "", "Number of Replicas(n)")
	optNumFaulty := getopt.Uint64('f', nFaulty, "", "Numer of Faulty Replicas(f) [Cannot exceed n-1/2]")
	optBlockSize := getopt.Uint64('b', blkSize, "", "Number of commands per block(b)")
	optDelay := getopt.Uint64('d', 10000, "", "Network Delay(d) [in milliseconds]")
	optBasePort := getopt.Uint32('p', basePort, "", "Base port for repicas. The nodes will use ports starting from this port number.")
	optClientBasePort := getopt.Uint32('c', clientBasePort, "", "Base port for clients. The clients will use these ports to talk to the nodes.")
	help := getopt.BoolLong("help", 'h', "Show this help message and exit")
	optOutDir := getopt.String('o', outDir, "Output Directory for the config files")

	getopt.Parse()
	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	if nReplicas == *optNumReplica {
		nFaulty = *optNumFaulty
	} else {
		nReplicas = *optNumReplica
		tempFaulty := uint64((nReplicas - 1) / 2)
		if *optNumFaulty < tempFaulty {
			nFaulty = *optNumFaulty
		} else {
			nFaulty = tempFaulty
		}
	}
	blkSize = *optBlockSize
	delta = float64(*optDelay) / 1000.0
	basePort = *optBasePort
	clientBasePort = *optClientBasePort
	outDir = *optOutDir

	// TODO: Pick PKI Algorithm from command line

	// Fetching the context
	alg := secp256k1.Secp256k1Context

	// NodeConfig
	nodeMap := make(map[uint64]*e2cconfig.NodeDataConfig)
	// Address Map for Protocol Nodes
	addrMap := make(map[uint64]*config.Address)
	// Address Map for Clients
	cliMap := make(map[uint64]*config.Address)
	// Public Key Map
	pubKeyMap := make(map[uint64][]byte)

	var err error

	for i := uint64(0); i < nReplicas; i++ {
		// Create a config
		nodeMap[i] = &e2cconfig.NodeDataConfig{}

		// Create Address and set it in the next loop
		addrMap[i] = &config.Address{}
		addrMap[i].IP = defaultIP
		addrMap[i].Port = fmt.Sprintf("%d", basePort+uint32(i))

		cliMap[i] = &config.Address{}
		cliMap[i].IP = defaultIP
		cliMap[i].Port = fmt.Sprintf("%d", clientBasePort+uint32(i))

		// Generate keypairs
		pvtKey, pubkey := alg.KeyGen()

		nodeMap[i].CryptoCon = &config.CryptoConfig{}
		nodeMap[i].CryptoCon.KeyType = alg.Type()
		nodeMap[i].CryptoCon.PvtKey, err = pvtKey.Raw()
		check(err)

		// Set it in the next loop
		pubKeyMap[i], err = pubkey.Raw()
		check(err)

		// Setup Protocol Configuration
		nodeMap[i].ProtConfig = &e2cconfig.E2CConfig{}
		nodeMap[i].ProtConfig.Id = i
		nodeMap[i].ProtConfig.Delta = delta
		nodeMap[i].ProtConfig.Info = &e2cconfig.ProtoInfo{}
		nodeMap[i].ProtConfig.Info.NodeSize = nReplicas
		nodeMap[i].ProtConfig.Info.Faults = nFaulty
		nodeMap[i].ProtConfig.Info.BlockSize = blkSize
	}

	for i := uint64(0); i < nReplicas; i++ {
		nodeMap[i].NetConfig = &config.NetConfig{}
		nodeMap[i].NetConfig.NodeAddressMap = addrMap

		nodeMap[i].ClientNetConfig = &config.NetConfig{}
		nodeMap[i].ClientNetConfig.NodeAddressMap = cliMap

		nodeMap[i].CryptoCon.NodeKeyMap = pubKeyMap
	}

	// Write Node Configs
	for i := uint64(0); i < nReplicas; i++ {
		// Open File
		fmt.Println("Processing Node:", i)
		fname := fmt.Sprintf(outDir+nodeConfigFile, i)
		// Serialize NodeConfig and Write to file
		nc := e2cconfig.NewNodeConfig(nodeMap[i])
		io.WriteToFile(nc, fname)
	}

	// Write a config for any client
	clientConfig := &e2cconfig.ClientDataConfig{}

	// Setup cryptographic configurations for the client
	clientConfig.CryptoCon = &config.CryptoConfig{}
	clientConfig.CryptoCon.KeyType = alg.Type()
	pvtKey, _ := alg.KeyGen()
	clientConfig.CryptoCon.PvtKey, err = pvtKey.Raw()
	check(err)
	clientConfig.CryptoCon.NodeKeyMap = pubKeyMap

	// Setup networking configurations for the client
	clientConfig.NetConfig = &config.NetConfig{}
	clientConfig.NetConfig.NodeAddressMap = cliMap

	// Setup Protocol Configurations
	clientConfig.Info = &e2cconfig.ProtoInfo{}
	clientConfig.Info.NodeSize = nReplicas
	clientConfig.Info.BlockSize = blkSize

	fname := fmt.Sprintf(outDir + clientFile)
	cc := e2cconfig.NewClientConfig(clientConfig)
	io.WriteToFile(cc, fname)
}
