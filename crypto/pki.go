package crypto

import (
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

var mapType = make(map[string]PKIAlgo)

// AddPKIAlgo A function that implements the interface must call this function to register their Algorithm
func AddPKIAlgo(name string, ctx PKIAlgo) {
	mapType[name] = ctx
}

// CheckExists checks if the algorithm is registered in the interface
func CheckExists(name string) bool {
	_, exists := mapType[name]
	return exists
}

// GetAlgo fetches the algorithm
// If the algorithm does not exist, it returns nil
func GetAlgo(name string) (alg PKIAlgo, exists bool) {
	alg, exists = mapType[name]
	return alg, exists
}

// PKIAlgo is the PKI Scheme Interface
type PKIAlgo interface {
	KeyGen() (PrivKey, PubKey)
	Type() string
	PubKeyFromBytes([]byte) PubKey
	PrivKeyFromBytes([]byte) PrivKey
}

// Key interface for cross compatibility with libp2p-crypto
type Key p2pcrypto.Key

// PubKey interface defines what a Public Key should look like
type PubKey p2pcrypto.PubKey

// PrivKey interface defines what a PrivKey should look like
type PrivKey p2pcrypto.PrivKey
