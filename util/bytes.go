package util

import (
	"encoding/hex"

	"github.com/adithyabhatkajake/libchatter/crypto"
)

// StringToByte converts a string into a byte array
func StringToByte(str string) []byte {
	b := []byte(str)
	return b
}

// HexStringToByte converts a hex string into bytes
func HexStringToByte(str string) []byte {
	h, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return h
}

// BytesToHexString converts byte arrays into a hex string */
func BytesToHexString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HashToString converts a hash into a hex string
func HashToString(nbytes [crypto.HashLen]byte) string {
	bytes := nbytes[:]
	return hex.EncodeToString(bytes)
}
