package e2c

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	pb "github.com/golang/protobuf/proto"
)

// Broadcast broadcasts a byte to all the nodes
func (e *E2C) Broadcast(m *msg.E2CMsg) error {
	e.netMutex.Lock()
	defer e.netMutex.Unlock()
	data, err := pb.Marshal(m)
	if err != nil {
		return err
	}
	for idx, s := range e.streamMap {
		_, err = s.Write(data)
		if err != nil {
			log.Error("Error while sending to node", idx)
			log.Error("Error:", err)
			return err
		}
		err = s.Flush()
		if err != nil {
			log.Error("Error while sending to node", idx)
			log.Error("Error:", err)
			return err
		}
	}
	return nil
}
