package synchs

import (
	"bufio"

	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
	pb "github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
)

// Implement how to talk to clients
const (
	ClientProtocolID = "synchs/client/0.0.1"
)

// ClientMsgHandler defines how to talk to client messages
func (shs *SyncHS) ClientMsgHandler(s network.Stream) {
	// A buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	shs.cliMutex.Lock()
	shs.cliMap[rw] = true
	shs.cliMutex.Unlock()
	// Event Handler
	for {
		// Receive a message from a client and process them
		len, err := rw.Read(buf)
		if err != nil {
			shs.errCh <- err
			log.Error("Error receiving a message from the client-", err)
			// Remove rw from cliMap after disconnection
			shs.cliMutex.Lock()
			_, exists := shs.cliMap[rw]
			// Delete only if it is in the map
			if exists {
				delete(shs.cliMap, rw)
			}
			shs.cliMutex.Unlock()
			return
		}
		// Send a copy for reacting
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React concurrently
		go shs.react(msgBuf)
	}
}

// ClientBroadcast sends a protocol message to all the clients known to this instance
func (shs *SyncHS) ClientBroadcast(m *msg.SyncHSMsg) {
	data, err := pb.Marshal(m)
	if err != nil {
		log.Error("Failed to send message", m, "to client")
		return
	}
	shs.cliMutex.Lock()
	defer shs.cliMutex.Unlock()
	for cliBuf := range shs.cliMap {
		log.Trace("Sending to", cliBuf)
		cliBuf.Write(data)
		cliBuf.Flush()
	}
	log.Trace("Finish client broadcast for", m)
}
