package e2c

var (
	// DefaultLeaderID is the ID of the Replica that the protocol starts with
	DefaultLeaderID uint64 = 1
)

func (e *E2C) changeLeader() {
	e.leader = (e.leader + 1) % uint64(len(e.pMap))
}
