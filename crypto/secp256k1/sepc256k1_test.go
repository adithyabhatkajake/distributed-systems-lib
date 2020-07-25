package secp256k1_test

import (
	"fmt"
	"testing"

	// To use require
	"github.com/stretchr/testify/require"

	"github.com/adithyabhatkajake/libe2c/crypto"
	"github.com/adithyabhatkajake/libe2c/crypto/secp256k1"
)

var (
	ctx     = secp256k1.Secp256k1Context
	msg     = []byte("Hello world")
	msgHash = crypto.DoHash(msg)
)

// Test if the key pair is generated properly and basic signatures are performed
// correctly
func TestKeyGen(t *testing.T) {
	// Generate a key pair and output it
	fmt.Println("Testing", ctx.Type())
	pvtkey, pubkey := ctx.KeyGen()
	require.NotEqual(t, nil, pvtkey)
	require.NotEqual(t, nil, pubkey)
}

func TestPubKeyCodec(t *testing.T) {
	_, pubkey := ctx.KeyGen()
	fmt.Println("Encoding Public Key")
	pubkeyBytes, err := pubkey.Raw()
	require.Equal(t, nil, err)
	pubkey_regen := ctx.PubKeyFromBytes(pubkeyBytes)
	fmt.Println("Decoding Public Key")
	newBytes, err2 := pubkey_regen.Raw()
	require.Equal(t, err2, nil)
	require.Equal(t, newBytes, pubkeyBytes)
	fmt.Println("Codec for Public Key passed")
}

func TestPrivKeyCodec(t *testing.T) {
	pvtkey, _ := ctx.KeyGen()
	fmt.Println("Encoding Private Key")
	privkeyBytes, err := pvtkey.Raw()
	require.Equal(t, err, nil)
	privkey_regen := ctx.PrivKeyFromBytes(privkeyBytes)
	fmt.Println("Decoding Private Key")
	newBytes, err2 := privkey_regen.Raw()
	require.Equal(t, nil, err2)
	require.Equal(t, privkeyBytes, newBytes)
	fmt.Println("Codec for Private Key passed")
}

func TestSignAndVerify(t *testing.T) {
	pvtkey, pubkey := ctx.KeyGen()
	sig, err := pvtkey.Sign(msg)
	if err != nil {
		panic(err)
	}
	status, err2 := pubkey.Verify(msg, sig)
	require.Equal(t, nil, err2)
	require.Equal(t, true, status)
}
