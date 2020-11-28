package model

import (
	"context"
	"fmt"
	"time"
)

type stateMachineImpl struct {
	tupleImpl
	sm           *StateMachineModel
	started      bool
	timeoutTimer *time.Timer
}

func (s *stateMachineImpl) GetStateTimeoutTimer() *time.Timer {
	return s.timeoutTimer
}

func (s *stateMachineImpl) SetStateTimeoutTimer(timer *time.Timer) {
	s.timeoutTimer = timer
}

func (s *stateMachineImpl) SetStarted(started bool) {
	s.started = started
}

func (s *stateMachineImpl) IsStarted() bool {
	return s.started
}

func NewStateMachineTuple(smc StateMachineModel, values map[string]interface{}) (StateMachineTuple, error) {
	valsNew := map[string]interface{}{}
	for k, v := range values {
		valsNew[k] = v
	}
	valsNew["sm_state"] = smc.InitialState
	tupleImplI, err := NewTuple(TupleType(smc.Descriptor.Name), valsNew)
	if err != nil {
		return nil, err
	}
	ti := tupleImplI.(*tupleImpl)
	smt := &stateMachineImpl{
		sm: &smc,
	}
	smt.tupleImpl.tuples = ti.tuples
	smt.tupleImpl.key = ti.GetKey()
	smt.tupleImpl.td = ti.GetTupleDescriptor()
	smt.tupleImpl.tupleType = ti.tupleType

	return smt, nil
}

func (s *stateMachineImpl) Start() {
	panic("implement me")
}

func (s *stateMachineImpl) SetState(ctx context.Context, state string) {
	s.SetString(ctx, "sm_state", state)
}

func (s *stateMachineImpl) GetState() string {
	if state, err := s.GetString("sm_state"); err == nil {
		return state
	}
	return ""
}

func (s *stateMachineImpl) GetStateMachine() *StateMachineModel {
	return s.sm
}

func (s *stateMachineImpl) CancelTimer() {
	panic("implement me")
}

func (s *stateMachineImpl) SetTimer() {
	panic("implement me")
}

//func (s *stateMachineImpl) Start() {
//	if s.isStarted {
//		return
//	}
//	s.currentState = s.sm.InitialState
//	s.timer = time.NewTimer(time.Duration(s.sm.States[s.currentState].Timeout) * time.Millisecond)
//	go func() {
//		<-s.timer.C
//		//todo: assert a time event!
//		//set timer keys, there will be a rule which matches and changes the state to timeout state
//	}()
//	s.isStarted = true
//}

func (s *stateMachineImpl) SetString(ctx context.Context, name string, value string) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}
func (s *stateMachineImpl) SetInt(ctx context.Context, name string, value int) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}
func (s *stateMachineImpl) SetLong(ctx context.Context, name string, value int64) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}
func (s *stateMachineImpl) SetDouble(ctx context.Context, name string, value float64) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}
func (s *stateMachineImpl) SetBool(ctx context.Context, name string, value bool) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}

func (s *stateMachineImpl) SetDatetime(ctx context.Context, name string, value time.Time) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}

func (s *stateMachineImpl) SetValue(ctx context.Context, name string, value interface{}) (err error) {
	return s.validateAndCallListener(ctx, name, value)
}

func (s *stateMachineImpl) validateAndCallListener(ctx context.Context, name string, value interface{}) (err error) {

	if s.isKeyProp(name) {
		return fmt.Errorf("Cannot change a key property [%s] for type [%s]", name, s.td.Name)
	}

	err = s.validateNameValue(name, value)
	if err != nil {
		return err
	}
	if s.tuples[name] != value {
		s.tuples[name] = value
		callChangeListener(ctx, s, name)
	}
	return nil
}
