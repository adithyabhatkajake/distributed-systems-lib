package crypto

import (
	"crypto/sha256"
)

const (
	// HashLen is the length of SHA256 output in bytes [32]
	HashLen = sha256.Size
)

// Hash is a fixed sized array type holding bytes of hash
type Hash [HashLen]byte

// DoHash takes bytes and outputs a 32 byte array
func DoHash(bytes []byte) Hash {
	h := sha256.Sum256(bytes)
	return h
}

// GetBytes returns []byte
func (h Hash) GetBytes() []byte {
	return h[:]
}

// ToHash converts []byte into [hashlen]byte
func ToHash(b []byte) Hash {
	var h Hash
	copy(h[:], b[0:HashLen])
	return h
}
