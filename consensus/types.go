package consensus

import (
	"github.com/adithyabhatkajake/libe2c/config"
	"github.com/adithyabhatkajake/libe2c/net"
)

/* This consensus algorithm just sends everybody a message and exits */

// Every consensus algorithm must implement three things.
// 1. Init(NodeDataConfig) - Prepare for consensus
// 2. Setup(channel) - Setup handlers and other stuff for protocol
// 3. Start() - Start the consensus

//Protocol is an interface to a generic protocol
type Protocol interface {
	Init(*config.NodeConfig)
	Start(*net.Network)
}
