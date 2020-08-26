package synchs

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
)

// QuitView quits the view
func (shs *SyncHS) QuitView() {
	log.Info("Quitting view ", shs.view)
	cert := &msg.Certificate{}
	var ids []uint64
	var sigs [][]byte // An array of byte arrays
	idx := uint64(0)
	shs.blLock.RLock()
	ids = make([]uint64, len(shs.blameMap[shs.view]))
	sigs = make([][]byte, len(shs.blameMap[shs.view]))
	for origin, bl := range shs.blameMap[shs.view] {
		ids[idx] = origin
		sigs[idx] = bl.Signature
	}
	shs.blLock.RUnlock()
	cert.Ids = ids
	cert.Signatures = sigs
	qv := &msg.QuitView{}
	qv.BlCert = cert
	m := &msg.SyncHSMsg{}
	m.Msg = &msg.SyncHSMsg_QV{QV: qv}
	shs.Broadcast(m)
	log.Debug("Finished Quitting view", shs.view)
}
