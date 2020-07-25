package main

import (
	"fmt"
	"os"

	"github.com/adithyabhatkajake/libe2c/io"

	"github.com/pborman/getopt"

	"github.com/adithyabhatkajake/libe2c/config"
	"github.com/adithyabhatkajake/libe2c/crypto/secp256k1"
)

var (
	nReplicas uint64  = 10
	nFaulty   uint64  = 4
	blkSize   uint64  = 1
	delta     float64 = 1 // in seconds
	basePort  uint32  = 10000
	outDir    string  = "testData/"
)

var (
	generalConfigFile = "nodes.txt"
	nodeConfigFile    = "nodes-%d.txt"
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
	optDelay := getopt.Uint64('d', 1000, "", "Network Delay(d) [in milliseconds]")
	optBasePort := getopt.Uint32('p', basePort, "", "Base port for repicas. The nodes will use ports starting from this port number.")
	help := getopt.BoolLong("help", 'h', "Show this help message and exit")
	optOutDir := getopt.String('o', outDir, "Output Directory for the config files")

	getopt.Parse()
	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	nReplicas = *optNumReplica
	nFaulty = *optNumFaulty
	blkSize = *optBlockSize
	delta = float64(*optDelay) / 1000.0
	basePort = *optBasePort
	outDir = *optOutDir

	// TODO: Pick PKI Algorithm from command line

	// Fetching the context
	alg := secp256k1.Secp256k1Context

	// NodeConfig
	nodeMap := make(map[uint64]*config.NodeDataConfig)
	// Address Map
	addrMap := make(map[uint64]*config.Address)
	// Public Key Map
	pubKeyMap := make(map[uint64][]byte)

	var err error

	for i := uint64(0); i < nReplicas; i++ {
		// Create a config
		nodeMap[i] = &config.NodeDataConfig{}

		// Create Address and set it in the next loop
		addrMap[i] = &config.Address{}
		addrMap[i].IP = "127.0.0.1"
		addrMap[i].Port = fmt.Sprintf("%d", basePort+uint32(i))

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
		nodeMap[i].ProtConfig = &config.E2CConfig{}
		nodeMap[i].ProtConfig.Id = i
		nodeMap[i].ProtConfig.NodeSize = nReplicas
		nodeMap[i].ProtConfig.BlockSize = blkSize
		nodeMap[i].ProtConfig.Delta = delta
	}

	for i := uint64(0); i < nReplicas; i++ {
		nodeMap[i].NetConfig = &config.NetConfig{}
		nodeMap[i].NetConfig.NodeAddressMap = addrMap
		nodeMap[i].CryptoCon.NodeKeyMap = pubKeyMap
	}

	for i := uint64(0); i < nReplicas; i++ {
		// Open File
		fmt.Println("Processing Node:", i)
		fname := fmt.Sprintf(outDir+nodeConfigFile, i)
		// Serialize NodeConfig and Write to file
		nc := config.NewNodeConfig(nodeMap[i])
		io.WriteToFile(nc, fname)
	}
}
