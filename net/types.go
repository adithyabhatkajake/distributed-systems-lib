package net

/* Networking module needs a config that implements the following interface */

// Conf implements networking configuration
type Conf interface {
	GetID() uint64
	GetP2PAddrFromID(uint64) string
}
