package e2c

import (
	"bufio"

	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	pb "github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
)

// Implement how to talk to clients
const (
	ClientProtocolID = "e2c/client/0.0.1"
)

// ClientMsgHandler defines how to talk to client messages
func (e *E2C) ClientMsgHandler(s network.Stream) {
	// A buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	e.cliMutex.Lock()
	e.cliMap[rw] = true
	e.cliMutex.Unlock()
	// Event Handler
	for {
		// Receive a message from a client and process them
		len, err := rw.Read(buf)
		if err != nil {
			e.errCh <- err
			log.Error("Error receiving a message from the client-", err)
			// Remove rw from cliMap after disconnection
			e.cliMutex.Lock()
			_, exists := e.cliMap[rw]
			// Delete only if it is in the map
			if exists {
				delete(e.cliMap, rw)
			}
			e.cliMutex.Unlock()
			return
		}
		// Send a copy for reacting
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React concurrently
		go e.react(msgBuf)
	}
}

// ClientBroadcast sends a protocol message to all the clients known to this instance
func (e *E2C) ClientBroadcast(m *msg.E2CMsg) {
	data, err := pb.Marshal(m)
	if err != nil {
		log.Error("Failed to send message", m, "to client")
		return
	}
	e.cliMutex.Lock()
	defer e.cliMutex.Unlock()
	for cliBuf := range e.cliMap {
		log.Trace("Sending to", cliBuf)
		cliBuf.Write(data)
		cliBuf.Flush()
	}
	log.Trace("Finish client broadcast for", m)
}
