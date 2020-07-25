package crypto

import (
	crand "crypto/rand"
)

// RandBytes is an abstraction over OS randomness
func RandBytes(numbytes int) []byte {
	rbytes := make([]byte, numbytes)
	_, err := crand.Read(rbytes)
	if err != nil {
		panic(err)
	}
	return rbytes
}
