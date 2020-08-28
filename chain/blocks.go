package chain

import (
	"bytes"

	"github.com/adithyabhatkajake/libchatter/crypto"
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

// GetHashBytes returns crypto.Hash from []bytes
func (b *Block) GetHashBytes() crypto.Hash {
	var x crypto.Hash
	copy(x[:], b.BlockHash)
	return x
}

// IsValid checks if the block is valid
func (b *Block) IsValid() bool {
	// Check if the hash is correctly computed
	if !bytes.Equal(b.GetHash().GetBytes(), b.BlockHash) {
		return false
	}
	return true
}
