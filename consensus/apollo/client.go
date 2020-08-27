package apollo

import (
	"bufio"

	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/apollo"
	pb "github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
)

// Implement how to talk to clients
const (
	ClientProtocolID = "apollo/client/0.0.1"
)

// ClientMsgHandler defines how to talk to client messages
func (n *Apollo) ClientMsgHandler(s network.Stream) {
	// A buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	n.cliMutex.Lock()
	n.cliMap[rw] = true
	n.cliMutex.Unlock()
	// Event Handler
	for {
		// Receive a message from a client and process them
		len, err := rw.Read(buf)
		if err != nil {
			n.errCh <- err
			log.Error("Error receiving a message from the client-", err)
			// Remove rw from cliMap after disconnection
			n.cliMutex.Lock()
			_, exists := n.cliMap[rw]
			// Delete only if it is in the map
			if exists {
				delete(n.cliMap, rw)
			}
			n.cliMutex.Unlock()
			return
		}
		// Send a copy for reacting
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React concurrently
		go n.react(msgBuf)
	}
}

// ClientBroadcast sends a protocol message to all the clients known to this instance
func (n *Apollo) ClientBroadcast(m *msg.ApolloMsg) {
	data, err := pb.Marshal(m)
	if err != nil {
		log.Error("Failed to send message", m, "to client")
		return
	}
	n.cliMutex.Lock()
	defer n.cliMutex.Unlock()
	for cliBuf := range n.cliMap {
		log.Trace("Sending to", cliBuf)
		cliBuf.Write(data)
		cliBuf.Flush()
	}
	log.Trace("Finish client broadcast for", m)
}
