package bench_test

import (
	"bufio"
	"context"
	"sync"
	"testing"

	"github.com/adithyabhatkajake/libchatter/crypto/secp256k1"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// BenchmarkSend benchmarks libp2p stream send cost
func testSend(b *testing.B) {
	// Generate two peers
	const testID = "test/testProtocol/0.0.1"
	ctx1 := context.Background()
	ctx2 := context.Background()
	alg := secp256k1.Secp256k1Context
	sk1, pk1 := alg.KeyGen()
	sk2, pk2 := alg.KeyGen()
	h1, err := libp2p.New(ctx1,
		libp2p.Identity(sk1),
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),
	)
	check(err)
	h2, err := libp2p.New(ctx2,
		libp2p.Identity(sk2),
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2001"),
	)
	check(err)
	h1.SetStreamHandler(testID, streamHandler)
	h2.SetStreamHandler(testID, streamHandler)
	p1, err := peer.IDFromPublicKey(pk1)
	check(err)
	p2, err := peer.IDFromPublicKey(pk2)
	check(err)
	p1addr := peer.AddrInfo{
		ID:    p1,
		Addrs: h1.Addrs(),
	}
	p2addr := peer.AddrInfo{
		ID:    p2,
		Addrs: h2.Addrs(),
	}
	wg := &sync.WaitGroup{}
	// Connect h1 to h2
	go func() {
		wg.Add(1)
		defer wg.Done()
		err := h1.Connect(ctx1, p2addr)
		for err != nil {
			err = h1.Connect(ctx1, p2addr)
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		// Connect h2 to h1
		err := h2.Connect(ctx2, p1addr)
		for err != nil {
			err = h2.Connect(ctx2, p1addr)
		}
	}()
	wg.Wait()
	h1s, err := h1.NewStream(ctx1, h2.ID(), testID)
	check(err)
	data := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h1s.Write(data)
	}
}

func streamHandler(s network.Stream) {
	r := bufio.NewReader(s)
	var buf []byte = make([]byte, 100)
	for {
		_, _ = r.Read(buf)
	}
}
