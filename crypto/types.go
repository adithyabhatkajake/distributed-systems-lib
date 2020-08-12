package crypto

import p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"

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

// Config defines what a cryptographic config should look like
type Config interface {
	GetMyKey() PrivKey
	GetPubKeyFromID(uint64) PubKey
}
