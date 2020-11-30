package model

import (
	"context"
	"time"
)

type StateMachineTuple interface {
	MutableTuple
	CancelTimer()
	SetTimer()
	Start()
	SetState(ctx context.Context, state string)
	GetState() string
	//SetStateMachine(sm *StateMachine)
	GetStateMachine() *StateMachine
	IsStarted() bool
	SetStarted(started bool)
	GetStateTimeoutTimer() *time.Timer
	SetStateTimeoutTimer(timer *time.Timer)
}
