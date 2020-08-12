package chain

import (
	"bytes"

	"github.com/adithyabhatkajake/libe2c/crypto"
	pb "github.com/golang/protobuf/proto"
)

// GetHash computes the hash from the block data
func (b *Block) GetHash() crypto.Hash {
	data, err := pb.Marshal(b.Data)
	if err != nil {
		panic(err)
	}
	return crypto.DoHash(data)
}

// IsValid checks if the block is valid
func (b *Block) IsValid() bool {
	// Check if the hash is correctly computed
	if !bytes.Equal(b.GetHash(), b.BlockHash) {
		return false
	}
	return true
}
