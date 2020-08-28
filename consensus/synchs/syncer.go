package synchs

// func (n *SyncHS) syncer() {
// 	candBlocks := make(map[uint64]*chain.Block)
// 	select {
// 	case prop, ok := <-n.syncObj.propChannel:
// 		log.Trace("Got a proposal for syncing")
// 	case candBlk, ok := <-n.syncObj.candBlockChannel:
// 		log.Trace("Got a candidate block for proposal")
// 		// // Do we have a certificate for the previous block?
// 		// n.certMapLock.RLock()
// 		// n.bc.ChainLock.RLock()
// 		// _, exists = n.certMap[n.bc.Head]
// 		// n.bc.ChainLock.RUnlock()
// 		// n.certMapLock.RUnlock()
// 		// if !exists {
// 		// 	log.Debug("No certificate for head found. Aborting proposal")
// 		// 	continue
// 		// }
// 	}
// }
