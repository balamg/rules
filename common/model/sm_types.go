package model

import (
	"context"
	"errors"
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
	SetPreviousState(state string)
	GetPreviousState() string
}

//state machine callbacks
type SmEntryActionFunction func(ctx context.Context, rs RuleSession, ruleName string,
	previousState string, currentState string,
	tuples map[TupleType]Tuple, ruleCtx RuleContext)

var (
	smEntryFns = make(map[string]SmEntryActionFunction)
)

// RegisterActionFunction registers the specified ActionFunction
func RegisterSmEntryFunction(id string, fn SmEntryActionFunction) error {

	if fn == nil {
		return errors.New("cannot register 'nil' SmEntryFunction")
	}

	if _, dup := smEntryFns[id]; dup {
		return errors.New("SmEntryFunction already registered: " + id)
	}

	smEntryFns[id] = fn

	return nil
}

// GetActionFunction gets specified ActionFunction
func GetSmEntryFunction(id string) SmEntryActionFunction {
	return smEntryFns[id]
}
