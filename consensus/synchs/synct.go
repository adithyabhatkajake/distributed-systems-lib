package synchs

import (
	"github.com/adithyabhatkajake/libe2c/chain"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
)

type syncT struct {
	propChannel      chan *msg.Proposal
	candBlockChannel chan *chain.Block
	certChannel      chan *msg.BlockCertificate
}

func newSyncT() syncT {
	s := syncT{}
	s.propChannel = make(chan *msg.Proposal, 10)
	s.candBlockChannel = make(chan *chain.Block, 10)
	s.certChannel = make(chan *msg.BlockCertificate, 10)
	return s
}
