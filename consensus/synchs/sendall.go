package synchs

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
	pb "github.com/golang/protobuf/proto"
)

// Broadcast broadcasts a byte to all the nodes
func (shs *SyncHS) Broadcast(m *msg.SyncHSMsg) error {
	shs.netMutex.Lock()
	defer shs.netMutex.Unlock()
	data, err := pb.Marshal(m)
	if err != nil {
		return err
	}
	// If we fail to send a message to someone, continue
	for idx, s := range shs.streamMap {
		_, err = s.Write(data)
		if err != nil {
			log.Error("Error while sending to node", idx)
			log.Error("Error:", err)
		}
		err = s.Flush()
		if err != nil {
			log.Error("Error while sending to node", idx)
			log.Error("Error:", err)
		}
	}
	return nil
}
