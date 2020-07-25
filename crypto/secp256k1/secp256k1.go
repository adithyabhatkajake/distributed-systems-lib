package secp256k1

const (
	algoname = "SECP256K1"
	// PrivKeySize is the size of the private key
	// L62 of secp256.go
	PrivKeySize = 32
)

var (
	// Secp256k1Context implements PKIAlgo interface
	Secp256k1Context = Secp256k1{}
)
