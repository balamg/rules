package model

import (
	"time"
)

type StateMachineImpl struct {
	tupleImpl
	sm           *StateMachineModel
	currentState string
	isStarted    bool
	timer        *time.Timer
}

func (s *StateMachineImpl) CancelTimer() {
	panic("implement me")
}

func (s *StateMachineImpl) SetTimer() {
	panic("implement me")
}

func (s *StateMachineImpl) Start() {
	if s.isStarted {
		return
	}
	s.currentState = s.sm.InitialState
	s.timer = time.NewTimer(time.Duration(s.sm.States[s.currentState].Timeout) * time.Millisecond)
	go func() {
		<-s.timer.C
		//todo: assert a time event!
		//set timer keys, there will be a rule which matches and changes the state to timeout state
	}()
	s.isStarted = true
}

type StateMachineModel struct {
	Name         string
	InitialState string
	States       map[string]*SmState
	EndState     string
}
type SmState struct {
	Name         string
	Timeout      int
	TimeoutState string
	Transitions  map[string]*SmTransition
}

type SmTransition struct {
	ExitAction  string
	NextState   string
	EntryAction string
}
