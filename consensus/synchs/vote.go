package synchs

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/synchs"
	pb "github.com/golang/protobuf/proto"
)

// Vote channel
func (n *SyncHS) voteHandler() {
	// Map leader to a map from sender to its vote
	voteMap := make(map[uint64]map[uint64]*msg.Vote)
	isCertified := make(map[uint64]bool)
	for {
		v, ok := <-n.voteChannel
		if !ok {
			log.Error("Vote channel error")
			continue
		}
		// Check if the vote is valid
		isValid := n.isVoteValid(v)
		if !isValid {
			log.Error("Invalid vote message")
			continue
		}
		// Check if this the first vote for this block height
		_, exists := voteMap[v.Data.Block.Data.Index]
		if !exists {
			voteMap[v.Data.Block.Data.Index] = make(map[uint64]*msg.Vote)
			isCertified[v.Data.Block.Data.Index] = false
		}
		isCert, exists := isCertified[v.Data.Block.Data.Index]
		if exists && isCert {
			// This vote for this block is already certified, ignore
			continue
		}
		_, exists = voteMap[v.Data.Block.Data.Index][v.Origin]
		if exists {
			log.Debug("Duplicate vote received")
			continue
		}
		voteMap[v.Data.Block.Data.Index][v.Origin] = v
		if uint64(len(voteMap[v.Data.Block.Data.Index])) <= n.config.GetNumberOfFaultyNodes() {
			log.Debug("Not enough votes to build a certificate")
			continue
		}
		log.Debug("Building a certificate")
		// Our certificate for height v.Data.Block.Data.Index is ready now
		cert := NewCert(voteMap[v.Data.Block.Data.Index])
		bcert := &msg.BlockCertificate{}
		bcert.BCert = cert
		bcert.Data = v.Data
		// Add this to certMaps
		n.certMapLock.Lock()
		n.certMap[v.Data.Block.Data.Index] = bcert
		n.certMapLock.Unlock()
		isCertified[v.Data.Block.Data.Index] = true
		log.Trace("Built a certificate", cert.String())
		if n.config.GetID() == n.leader {
			go n.propose() // We propose here also hoping that there are sufficient comamnds to build a block
		}
	}
}

func (n *SyncHS) isVoteValid(v *msg.Vote) bool {
	data, err := pb.Marshal(v.Data)
	if err != nil {
		log.Error("Error marshalling vote data")
		log.Error(err)
		return false
	}
	sigOk, err := n.config.GetPubKeyFromID(v.Origin).Verify(data, v.Signature)
	if err != nil {
		log.Error("Vote Signature check error")
		log.Error(err)
		return false
	}
	if !sigOk {
		log.Error("Invalid vote Signature")
		return sigOk
	}
	data, err = pb.Marshal(v.Data.Block.Data)
	if err != nil {
		log.Error("Error marshalling block data from vote")
		log.Error(err)
		return false
	}
	sigOk, err = n.config.GetPubKeyFromID(v.Data.Block.Proposer).Verify(data, v.Data.Block.Signature)
	if err != nil {
		log.Error("Error while checking Block Signature in the vote")
		log.Error(err)
		return false
	}
	if !sigOk {
		log.Error("Invalid block Signature in vote")
		return sigOk
	}
	return sigOk
}
