package model

import "context"

type StateMachineTuple interface {
	MutableTuple
	CancelTimer()
	SetTimer()
	Start()
	SetState(ctx context.Context, state string)
	GetState() string
	//SetStateMachine(sm *StateMachineModel)
	GetStateMachine() *StateMachineModel
}
