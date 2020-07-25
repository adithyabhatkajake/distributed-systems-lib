package config

import (
	"fmt"
)

// GetP2PAddr returns a Libp2p compatible address from address strings
func (x *Address) GetP2PAddr() string {
	addr := fmt.Sprintf("/ip4/%s/tcp/%s", x.GetIP(), x.GetPort())
	return addr
}
