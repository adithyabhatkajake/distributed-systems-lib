package chain

import (
	"github.com/adithyabhatkajake/libe2c/crypto"
	pb "github.com/golang/protobuf/proto"
)

// GetHash returns the hash of the command
func (c *Command) GetHash() crypto.Hash {
	data, err := pb.Marshal(c)
	if err != nil {
		panic(err)
	}
	return crypto.DoHash(data)
}
