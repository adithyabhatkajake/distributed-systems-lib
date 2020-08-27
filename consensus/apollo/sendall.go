package apollo

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/apollo"
	pb "github.com/golang/protobuf/proto"
)

// Broadcast broadcasts the message to all the nodes
func (n *Apollo) Broadcast(m *msg.ApolloMsg) error {
	n.netMutex.Lock()
	defer n.netMutex.Unlock()
	data, err := pb.Marshal(m)
	if err != nil {
		return err
	}
	// If we fail to send a message to someone, continue
	for idx, s := range n.streamMap {
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

// Send sends a protocol message m to a particular node whose ID is receiverID
func (n *Apollo) Send(m *msg.ApolloMsg, receiverID uint64) error {
	data, err := pb.Marshal(m)
	if err != nil {
		return err
	}
	n.netMutex.Lock()
	defer n.netMutex.Unlock()
	s := n.streamMap[receiverID]
	_, err = s.Write(data)
	if err != nil {
		log.Error("Error while sending to node", receiverID)
		log.Error("Error:", err)
		return err
	}
	err = s.Flush()
	if err != nil {
		log.Error("Error while flushing the data to node", receiverID)
		log.Error("Error:", err)
		return err
	}
	return err
}
