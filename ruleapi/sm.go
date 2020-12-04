package ruleapi

import (
	"context"
	"fmt"
	"time"

	"github.com/project-flogo/rules/common/model"
)

type ActionCtx struct {
	sms        *model.StateMachines
	sm         *model.StateMachine
	state      *model.SmState
	transition *model.SmTransition
}

type StateTimeoutActionCtx struct {
	ActionCtx
}

type StateChangeActionCtx struct {
	ActionCtx
}

type StartChildSmActionCtx struct {
	ActionCtx
}

type ExitChildSmActionCtx struct {
	ActionCtx
}

func CreateRulesForSm(sms *model.StateMachines, sm *model.StateMachine) ([]model.Rule, error) {
	var rules []model.Rule

	for i := range sm.States {
		state := &sm.States[i]
		rls, err := CreateRulesForState(sms, sm, state)
		if err != nil {
			return nil, err
		}
		for i := range rls {
			rules = append(rules, rls[i])
		}
	}

	return rules, nil
}

func CreateRulesForState(sms *model.StateMachines, sm *model.StateMachine, state *model.SmState) ([]model.Rule, error) {
	smName := sm.Descriptor.Name
	var rules []model.Rule

	for i := range state.Transitions {

		transition := &state.Transitions[i]

		condition := transition.Condition
		ruleName := fmt.Sprintf("%s_%s_%s_%s", smName, state.State, transition.ToState, condition)

		rule := NewRule(ruleName)
		err := addDefinedCondition(rule, condition)
		if err != nil {
			return nil, err
		}

		err = addCurrentStateCondition(rule, sm, state)
		if err != nil {
			return nil, err
		}

		err = addParentSmConditions(rule, sms, sm)
		if err != nil {
			return nil, err
		}

		if transition.StartSm == "" {
			//its a regular state change
			sa := &StateChangeActionCtx{ActionCtx{
				sms:        sms,
				sm:         sm,
				state:      state,
				transition: transition,
			}}
			rule.SetAction(sa.stateChangeAction)
			rules = append(rules, rule)
		} else {
			//set action to startsm action
			startChildSmActionCtx := &StartChildSmActionCtx{ActionCtx{
				sms:        sms,
				sm:         sm,
				state:      state,
				transition: transition,
			}}
			rule.SetAction(startChildSmActionCtx.startChildSmAction)
			rules = append(rules, rule)

			//additional
			rule, err = defineExitChildSmRule(sms, sm, state, transition)
			if err != nil {
				return rules, err
			}
			rules = append(rules, rule)
		}

	}
	if state.Timeout > 0 {
		rule, err := defineStateTimeoutRule(sms, sm, state, nil)
		if err != nil {
			return rules, err
		}
		err = addParentSmConditions(rule, sms, sm)
		if err != nil {
			return rules, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func addCurrentStateCondition(rule model.MutableRule, sm *model.StateMachine,
	state *model.SmState) error {
	currStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", sm.Descriptor.Name, state.State)
	err := rule.AddExprCondition(currStateCondition, currStateCondition, nil)
	if err != nil {
		return err
	}
	return nil
}

func addDefinedCondition(rule model.MutableRule, condition string) error {
	err := rule.AddExprCondition(condition, condition, nil)
	if err != nil {
		return err
	}
	return nil
}

func addParentSmConditions(rule model.MutableRule, sms *model.StateMachines, sm *model.StateMachine) error {
	currentParent := sm.ParentSm
	currentParentState := sm.ParentState

	for currentParent != "" {
		matchParentStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", currentParent, currentParentState)
		err := rule.AddExprCondition(matchParentStateCondition, matchParentStateCondition, nil)
		if err != nil {
			return err
		}

		matchParentKeyCondition := fmt.Sprintf("$.%s.sm_key == $.%s.sm_key",
			currentParent, sm.Descriptor.Name)
		err = rule.AddExprCondition(matchParentKeyCondition, matchParentKeyCondition, nil)
		if err != nil {
			return err
		}

		parentSm := sms.GetSm(currentParent)
		if parentSm == nil {
			currentParent = ""
		} else {
			currentParent = parentSm.ParentSm
			currentParentState = parentSm.ParentState
		}
	}
	return nil
}

func defineExitChildSmRule(sms *model.StateMachines, sm *model.StateMachine,
	state *model.SmState, transition *model.SmTransition) (model.MutableRule, error) {
	ruleName := fmt.Sprintf("%s_%s_%s_exit", sm.Descriptor.Name,
		state.State, transition.StartSm)
	rule := NewRule(ruleName)

	err := addCurrentStateCondition(rule, sm, state)
	if err != nil {
		return nil, err
	}
	childSm := sms.GetSm(transition.StartSm)
	childEndStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'",
		childSm.Descriptor.Name, childSm.EndState)
	err = rule.AddExprCondition(childEndStateCondition, childEndStateCondition, nil)
	if err != nil {
		return nil, err
	}

	//equal keys rule
	equalKeysCondition := fmt.Sprintf("$.%s.sm_key == $.%s.sm_key", sm.Descriptor.Name,
		childSm.Descriptor.Name)
	err = rule.AddExprCondition(equalKeysCondition, equalKeysCondition, nil)
	if err != nil {
		return nil, err
	}
	exitChildSmActionCtx := &ExitChildSmActionCtx{ActionCtx{
		sms:        sms,
		sm:         sm,
		state:      state,
		transition: transition,
	}}
	err = addParentSmConditions(rule, sms, sm)
	if err != nil {
		return nil, err
	}
	rule.SetAction(exitChildSmActionCtx.exitChildSmAction)
	return rule, nil
}

func defineStateTimeoutRule(sms *model.StateMachines, sm *model.StateMachine,
	state *model.SmState, transition *model.SmTransition) (model.MutableRule, error) {

	ruleName := fmt.Sprintf("%s_%s_timeout", sm.Descriptor.Name, state.State)
	rule := NewRule(ruleName)

	err := addCurrentStateCondition(rule, sm, state)
	if err != nil {
		return nil, err
	}

	//$.sm1.sm_key == $.timer.ctx
	matchKeyExpr := fmt.Sprintf("$.%s.sm_key == $.timer.ctx", sm.Descriptor.Name)
	err = rule.AddExprCondition(matchKeyExpr, matchKeyExpr, nil)
	if err != nil {
		return nil, err
	}

	//$.sm1.sm_key == $.timer.ctx
	matchTimerRuleName := fmt.Sprintf("$.timer.ruleName == '%s'", ruleName)
	err = rule.AddExprCondition(matchTimerRuleName, matchTimerRuleName, nil)
	if err != nil {
		return nil, err
	}

	smt := &StateTimeoutActionCtx{ActionCtx{
		sms:        sms,
		sm:         sm,
		state:      state,
		transition: transition,
	}}
	rule.SetAction(smt.stateTimeoutAction)
	return rule, nil
}

func (a *StateChangeActionCtx) stateChangeAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	smName := a.sm.Descriptor.Name
	smTuple := tuples[model.TupleType(smName)]
	smt, _ := smTuple.(model.StateMachineTuple)

	fmt.Printf("state machine [%s]: state [%s] changed to state [%s]\n",
		smt.GetKey().String(), smt.GetState(), a.transition.ToState)
	smt.SetState(ctx, a.transition.ToState)
	_ = startTimeoutForCurrentState(smt, a.sm, rs)
}

func (a *StateTimeoutActionCtx) stateTimeoutAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	smName := a.sm.Descriptor.Name
	smTuple := tuples[model.TupleType(smName)]
	smt, _ := smTuple.(model.StateMachineTuple)

	fmt.Printf("state machine [%s]: state [%s] timed out, setting timeout state [%s]\n",
		smt.GetKey().String(), smt.GetState(), a.state.TimeoutState)
	smt.SetState(ctx, a.state.TimeoutState)
	_ = startTimeoutForCurrentState(smt, a.sm, rs)
}

func StartSm(ctx context.Context, rs model.RuleSession, sm *model.StateMachine, s model.StateMachineTuple) error {
	if s.IsStarted() {
		return nil
	}
	_ = startTimeoutForCurrentState(s, sm, rs)

	err := rs.Assert(ctx, s)
	if err != nil {
		return err
	}
	s.SetStarted(true)
	return nil
}

//to be called right after changing state to nextState
func startTimeoutForCurrentState(s model.StateMachineTuple, sm *model.StateMachine,
	rs model.RuleSession) error {
	timer := s.GetStateTimeoutTimer()
	if timer != nil {
		timer.Stop()
	}
	currentState := sm.GetSmForState(s.GetState())
	if currentState == nil {
		fmt.Printf("WARN: state machine [%s]: state not found [%s]\n",
			sm.Descriptor.Name, s.GetState())
		return nil
	}
	timeout := currentState.Timeout
	if timeout > 0 {
		timer = time.NewTimer(time.Duration(timeout) * time.Second)
		s.SetStateTimeoutTimer(timer)

		go func() {
			<-timer.C
			smKey := s.GetMap()["sm_key"].(string)
			ruleName := fmt.Sprintf("%s_%s_timeout", sm.Descriptor.Name, s.GetState())
			assertTimerTuple(rs, ruleName, smKey)
		}()
	}
	return nil
}

func assertTimerTuple(rs model.RuleSession, ruleName string, smKey string) {
	now := time.Now().UnixNano()

	timer, _ := model.NewTupleWithKeyValues("timer", now)
	_ = timer.SetString(context.TODO(), "ctx", smKey)
	_ = timer.SetString(context.TODO(), "ruleName", ruleName)

	_ = rs.Assert(context.TODO(), timer)

}
func (a *ExitChildSmActionCtx) exitChildSmAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	smName := a.sm.Descriptor.Name
	smTuple := tuples[model.TupleType(smName)]
	smt, _ := smTuple.(model.StateMachineTuple)

	fmt.Printf("exit child state machine [%s]: state [%s] changed to state [%s]\n",
		smt.GetKey().String(), smt.GetState(), a.transition.ToState)
	smt.SetState(ctx, a.transition.ToState)
	_ = startTimeoutForCurrentState(smt, a.sm, rs)
}

func (a *StartChildSmActionCtx) startChildSmAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {

	smName := a.sm.Descriptor.Name
	smTuple := tuples[model.TupleType(smName)]
	smt, _ := smTuple.(model.StateMachineTuple)

	smKey := smt.GetMap()["sm_key"]
	valMap := map[string]interface{}{"sm_key": smKey}
	childSm := a.sms.GetSm(a.transition.StartSm)

	smt, err := model.NewStateMachineTuple(childSm, valMap)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	err = StartSm(ctx, rs, childSm, smt)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
