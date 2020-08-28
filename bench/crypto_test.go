package bench_test

import (
	"testing"

	"github.com/adithyabhatkajake/libchatter/crypto/secp256k1"
)

func BenchmarkSign(b *testing.B) {
	alg := secp256k1.Secp256k1Context
	sk1, _ := alg.KeyGen()
	data := make([]byte, 1000)
	for i := 0; i < b.N; i++ {
		_, _ = sk1.Sign(data)
	}
}

func BenchmarkVerify(b *testing.B) {
	alg := secp256k1.Secp256k1Context
	sk1, pk1 := alg.KeyGen()
	data := make([]byte, 1000)
	sig, _ := sk1.Sign(data)
	for i := 0; i < b.N; i++ {
		_, _ = pk1.Verify(data, sig)
	}
}
