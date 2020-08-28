package io

import (
	"github.com/adithyabhatkajake/libchatter/crypto"
	secp256k1 "github.com/adithyabhatkajake/libchatter/crypto/secp256k1"
)

// regsiter all contexts
func init() {
	crypto.AddPKIAlgo(
		secp256k1.Secp256k1Context.Type(),
		secp256k1.Secp256k1Context)
}
