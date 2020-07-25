package secp256k1

import (
	"crypto/rand"

	"github.com/adithyabhatkajake/libe2c/crypto"
	underlyingSecp256k1 "github.com/libp2p/go-libp2p-core/crypto"
)

// Secp256k1 implements the PKIAlgo interface
type Secp256k1 struct{}

// Type returns the name of the algorithm
func (s Secp256k1) Type() string {
	return algoname
}

// KeyGen generates a pair of public and private key pair
func (s Secp256k1) KeyGen() (crypto.PrivKey, crypto.PubKey) {
	pvtkey, pubkey, err := underlyingSecp256k1.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		panic(err)
	}
	return pvtkey, pubkey
}

// SecpPubKey implements libp2p public key
type SecpPubKey underlyingSecp256k1.Secp256k1PublicKey

// SecpPrivKey implements libp2p private key
type SecpPrivKey underlyingSecp256k1.Secp256k1PrivateKey

// PrivKeyFromBytes unmarshals the private key from raw bytes
func (s Secp256k1) PrivKeyFromBytes(data []byte) crypto.PrivKey {
	pvtkey, err := underlyingSecp256k1.UnmarshalSecp256k1PrivateKey(data)
	if err != nil {
		panic(err)
	}
	return pvtkey
}

// PubKeyFromBytes unmarshals a public key from raw bytes
func (s Secp256k1) PubKeyFromBytes(data []byte) crypto.PubKey {
	pubkey, err := underlyingSecp256k1.UnmarshalSecp256k1PublicKey(data)
	if err != nil {
		panic(err)
	}
	return pubkey
}
