package bench_test

import (
	"testing"

	"github.com/adithyabhatkajake/libchatter/chain"
	pb "github.com/golang/protobuf/proto"
)

func BenchmarkCommandSerialize(b *testing.B) {
	c := &chain.Command{}
	c.Cmd = []byte("Test string to test parsing times")
	for i := 0; i < b.N; i++ {
		pb.Marshal(c)
	}
}

func BenchmarkCommandDeserialize(b *testing.B) {
	c := &chain.Command{}
	c2 := &chain.Command{}
	c.Cmd = []byte("Test string to test parsing times")
	data, err := pb.Marshal(c)
	check(err)
	for i := 0; i < b.N; i++ {
		pb.Unmarshal(data, c2)
	}
}
