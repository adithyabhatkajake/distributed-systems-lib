package e2c

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
)

// QuitView quits the view
func (e *E2C) QuitView() {
	log.Info("Quitting view ", e.view)
	cert := &msg.Certificate{}
	var ids []uint64
	var sigs [][]byte // An array of byte arrays
	idx := uint64(0)
	e.blLock.RLock()
	ids = make([]uint64, len(e.blameMap[e.view]))
	sigs = make([][]byte, len(e.blameMap[e.view]))
	for origin, bl := range e.blameMap[e.view] {
		ids[idx] = origin
		sigs[idx] = bl.Signature
	}
	e.blLock.RUnlock()
	cert.Ids = ids
	cert.Signatures = sigs
	qv := &msg.QuitView{}
	qv.BlCert = cert
	m := &msg.E2CMsg{}
	m.Msg = &msg.E2CMsg_QV{QV: qv}
	e.Broadcast(m)
	log.Debug("Finished Quitting view", e.view)
}
