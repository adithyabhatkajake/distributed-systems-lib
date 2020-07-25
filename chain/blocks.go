package chain

import (
	"sync"

	"github.com/adithyabhatkajake/libe2c/crypto"
)

// Block Data Structure
type Block struct {
	Index     uint64
	BlockHash crypto.Hash
	PrevHash  crypto.Hash
	Decision  bool
	Proposer  uint64
	cmds      []crypto.Hash // Array of Hashes
}

type void struct{}

// BlockChain is what we call a blockchain
type BlockChain struct {
	Chain map[[crypto.HashLen]byte]Block
	// A lock that we use to safely update the chain
	ChainLock sync.Mutex
	// A height block map
	HeightBlockMap map[uint64]Block
	// Unconfirmed Blocks
	UnconfirmedBlocks map[crypto.Hash]void
}
