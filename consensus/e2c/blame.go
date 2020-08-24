package e2c

import (
	"github.com/adithyabhatkajake/libe2c/log"
	msg "github.com/adithyabhatkajake/libe2c/msg/e2c"
	"github.com/adithyabhatkajake/libe2c/util"
	pb "github.com/golang/protobuf/proto"
)

func (e *E2C) startBlameTimer() {
	log.Info("Starting Blame Timer")
	e.blTimerLock.Lock()
	defer e.blTimerLock.Unlock()
	if e.blTimer != nil {
		e.blTimer.Reset()
		log.Debug("Blame timer reset")
		return
	}
	e.blTimer = util.NewTimer(func() {
		e.sendNPBlame()
	})
	e.blTimer.SetTime(e.config.GetNPBlameWaitTime())
	e.blTimer.Start()
	log.Debug("Finished starting a blame timer")
}

func (e *E2C) stopBlameTimer() {
	log.Info("Stopping Blame Timer")
	e.blTimerLock.Lock()
	defer e.blTimerLock.Unlock()
	if e.blTimer != nil {
		e.blTimer.Cancel()
		e.blTimer = nil
	}
}

func (e *E2C) resetBlameTimer() {
	log.Info("Resetting Blame Timer")
	e.blTimerLock.Lock()
	e.blTimer.Reset() // Reset the timer
	e.blTimerLock.Unlock()
}

func (e *E2C) sendNPBlame() {
	log.Warn("Sending an NP-Blame message")
	blame := &msg.NoProgressBlame{}
	blame.Blame = &msg.Blame{}
	blame.Blame.BlData = &msg.BlameData{}
	blame.Blame.BlData.BlameTarget = e.leader
	blame.Blame.BlData.View = e.view
	blame.Blame.BlOrigin = e.config.GetID()
	data, err := pb.Marshal(blame.Blame.BlData)
	if err != nil {
		log.Errorln("Error marshalling blame message", err)
		return
	}
	blame.Blame.Signature, err = e.config.GetMyKey().Sign(data)
	if err != nil {
		log.Errorln("Error Signing the blame message", err)
	}
	blMsg := &msg.E2CMsg{}
	blMsg.Msg = &msg.E2CMsg_Npblame{Npblame: blame}
	e.Broadcast(blMsg)
}

func (e *E2C) handleNoProgressBlame(bl *msg.NoProgressBlame) {
	log.Trace("Received a No-Progress blame against ",
		bl.Blame.BlData.BlameTarget, " from ", bl.Blame.BlOrigin)
	// Check if the blame is correct
	isValid := e.isNPBlameValid(bl)
	if !isValid {
		log.Debugln("Received an invalid blame message", bl.String())
		return
	}
	// Add it to the blame map
	// TODO
	// Check if the blame map has sufficient blames
	// If there are more than f blames, then initiate quit view
}

func (e *E2C) isNPBlameValid(bl *msg.NoProgressBlame) bool {
	log.Traceln("Function isNPBlameValid with input", bl.String())
	// Check if the blame is for the current leader
	if bl.Blame.BlData.BlameTarget != e.leader {
		log.Debug("Invalid Blame Target. Found", bl.Blame.BlData.BlameTarget,
			",Expected:", e.leader)
		return false
	}
	// Check if the view is correct!
	if bl.Blame.BlData.View != e.view {
		log.Debug("Invalid Blame View. Found", bl.Blame.BlData.View,
			",Expected:", e.view)
		return false
	}
	// Get bl data
	data, err := pb.Marshal(bl.Blame.BlData)
	if err != nil {
		log.Debug("Error Marshalling blame message")
		return false
	}
	// Check if the signature is correct
	isSigValid, err := e.config.GetPubKeyFromID(
		bl.Blame.BlOrigin).Verify(data, bl.Blame.Signature)
	if !isSigValid || err != nil {
		log.Debug("Invalid signature for blame message")
		return false
	}
	return true
}
