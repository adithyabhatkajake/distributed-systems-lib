package config

// Config defines a generic config
type Config interface {
	GetNumNodes() uint64
}
