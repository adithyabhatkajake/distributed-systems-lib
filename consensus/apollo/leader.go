package apollo

var (
	// DefaultLeaderID is the ID of the Replica that the protocol starts with
	DefaultLeaderID uint64 = 1
)

func (n *Apollo) changeLeader() {
	n.leaderLock.Lock()
	defer n.leaderLock.Unlock()
	// For now rotate leaders.
	// TODO - Add a proposer set, change based on the next element in the proposer set
	n.leader = (n.leader + 1) % n.config.GetNumNodes()
}

func (n *Apollo) nextLeader() uint64 {
	n.leaderLock.RLock()
	defer n.leaderLock.RUnlock()
	return (n.leader + 1) % n.config.GetNumNodes()
}

func (n *Apollo) currentLeader() uint64 {
	n.leaderLock.RLock()
	defer n.leaderLock.RUnlock()
	return n.leader
}

func (n *Apollo) senderForBlockIndex(bIdx uint64) uint64 {
	// We have received block 1, the sender for this should be node whose ID is 1
	N := n.config.GetNumNodes()
	return bIdx % N
}
