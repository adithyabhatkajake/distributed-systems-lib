package io

import (
	"github.com/adithyabhatkajake/libe2c/crypto"
	secp256k1 "github.com/adithyabhatkajake/libe2c/crypto/secp256k1"
)

// regsiter all contexts
func init() {
	crypto.AddPKIAlgo(
		secp256k1.Secp256k1Context.Type(),
		secp256k1.Secp256k1Context)
}
