package crypto_test

import (
	"testing"
	// To use require
	"github.com/stretchr/testify/require"

	"github.com/adithyabhatkajake/libe2c/crypto"
)

/* This test checks for the following:*/
/* 1. Bytes generated are unique */
/* 2. The functions do not panic */

func TestRandomNoError(t *testing.T) {
	_ = crypto.RandBytes(100)
	require.Equal(t, true, true)
}

func TestRandomIsRandom(t *testing.T) {
	randBytes := crypto.RandBytes(100)
	randBytesSecond := crypto.RandBytes(100)

	require.NotEqual(t, randBytes, randBytesSecond)
}
