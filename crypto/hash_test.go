package crypto_test

import (
	"testing"

	// To use require
	"github.com/stretchr/testify/require"

	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/util"
)

const (
	str = "Libe2c is awesome"
	// This is obtained by running `echo -n $str | sha256sum`
	expStr = "a50e9835e4ef523c7d2e19da376fe02ab75f3686ea1e477a88e1f4ee9a5a658a"
)

// Test if the function computes the hash correctly
// Test against the standard linux implementation in sha256
func TestHash(t *testing.T) {
	testStringInBytes := util.StringToByte(str)
	computedHash := crypto.DoHash(testStringInBytes)
	expectedHashBytes := util.HexStringToByte(expStr) // []byte
	var expectedHash crypto.Hash = make([]byte, crypto.HashLen)
	copy(expectedHash[:], expectedHashBytes) // Now it is [32]byte
	require.Equal(t, len(computedHash), len(expectedHash))
	require.Equal(t, len(computedHash), crypto.HashLen)
	require.Equal(t, expectedHash, computedHash)
}
