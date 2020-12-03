package ruleapi

import (
	"context"
	"fmt"
	"time"

	"github.com/project-flogo/rules/common/model"
)

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
		ruleName := fmt.Sprintf("%s_%s_%s_%s", smName, state.State, condition, transition.ToState)
		rule := NewRule(ruleName)

		currStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", smName, state.State)
		err := rule.AddExprCondition(currStateCondition, currStateCondition, &transition)
		if err != nil {
			return rules, err
		}
		err = rule.AddExprCondition(condition, condition, nil)
		if err != nil {
			return rules, err
		}
		err = addParentSmConditions(sms, sm, rule)
		if err != nil {
			return rules, err
		}

		if transition.StartSm == "" {
			//its a regular state change
			sa := SmActionContext{name: smName, smTrans: transition}
			rule.SetAction(sa.setSmTransitionAction)
			rules = append(rules, rule)
		} else {
			rule, err = defineChildStateStartRule(sms, sm, state, transition)
			if err != nil {
				return rules, err
			}
			rules = append(rules, rule)

			rule, err = defineChildStateExitRule(sms, sm, state, transition)
			if err != nil {
				return rules, err
			}
			rules = append(rules, rule)
		}

	}
	if state.Timeout > 0 {
		timeoutRule, err := setTimeoutRuleForState(smName, state)
		if err != nil {
			return rules, err
		}
		err = addParentSmConditions(sms, sm, timeoutRule)
		if err != nil {
			return rules, err
		}
		rules = append(rules, timeoutRule)
	}
	return rules, nil
}
func defineChildStateStartRule(sms *model.StateMachines, sm *model.StateMachine,
	state *model.SmState, transition *model.SmTransition) (model.MutableRule, error) {
	ruleName := fmt.Sprintf("%s_%s_%s_enter", sm.Descriptor.Name,
		state.State, transition.StartSm)
	rule := NewRule(ruleName)

	currStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", sm.Descriptor.Name, state.State)
	err := rule.AddExprCondition(currStateCondition, currStateCondition, &transition)
	if err != nil {
		return nil, err
	}

	//child sm is involved, so create a start rule and an exit rule
	startChildSmActionCtx := &StartChildSmActionContext{
		name:    sm.Descriptor.Name,
		sms:     sms,
		state:   state,
		smTrans: transition,
	}
	rule.SetAction(startChildSmActionCtx.startChildSmAction)

	return rule, nil
}
func defineChildStateExitRule(sms *model.StateMachines, sm *model.StateMachine,
	state *model.SmState, transition *model.SmTransition) (model.MutableRule, error) {
	ruleName := fmt.Sprintf("%s_%s_%s_exit", sm.Descriptor.Name,
		state.State, transition.StartSm)
	rule := NewRule(ruleName)

	currStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", sm.Descriptor.Name, state.State)
	err := rule.AddExprCondition(currStateCondition, currStateCondition, &transition)
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
	childSmExitCtx := &ChildSmExitActionContext{
		name:    sm.Descriptor.Name,
		sms:     sms,
		state:   state,
		smTrans: transition,
	}
	rule.SetAction(childSmExitCtx.childSmExitAction)
	return rule, nil
}

func addParentSmConditions(sms *model.StateMachines, sm *model.StateMachine, rule model.MutableRule) error {
	//add parent rules
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

func setTimeoutRuleForState(smName string, state *model.SmState) (model.MutableRule, error) {
	ruleName := fmt.Sprintf("%s_%s_timeout", smName, state.State)
	rule := NewRule(ruleName)

	//$.sm1.sm_state == 's1'
	currStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", smName, state.State)
	err := rule.AddExprCondition(currStateCondition, currStateCondition, nil)
	if err != nil {
		return nil, err
	}

	//$.sm1.sm_key == $.timer.ctx
	matchKeyExpr := fmt.Sprintf("$.%s.sm_key == $.timer.ctx", smName)
	err = rule.AddExprCondition(matchKeyExpr, matchKeyExpr, nil)
	if err != nil {
		return nil, err
	}

	smt := &SmTimeoutActionContext{smName, state}
	rule.SetAction(smt.TimeoutAction)
	return rule, nil
}

func (smCtx *SmTimeoutActionContext) TimeoutAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	smTuple := tuples[model.TupleType(smCtx.name)]
	if smTuple == nil {
		fmt.Printf("sm not found %s", smCtx.name)
		return
	}

	smt, ok := smTuple.(model.StateMachineTuple)
	if !ok {
		fmt.Printf("%s not of type statemachinetuple", smCtx.name)
		return
	}

	fmt.Printf("setting sm[%s] to next state [%s] from state [%s]\n",
		smt.GetKey().String(), smCtx.sm.TimeoutState, smt.GetState())
	smt.SetState(ctx, smCtx.sm.TimeoutState)
	timerEvent := tuples[model.TupleType("timer")]
	if timerEvent == nil {
		fmt.Printf("timer event not found\n")
		return
	}
	rs.Delete(ctx, timerEvent)

	_ = startTimeoutForCurrentState(smt, rs)
}

type SmTimeoutActionContext struct {
	name string
	sm   *model.SmState
}
type SmActionContext struct {
	name    string
	smTrans *model.SmTransition
}
type StartChildSmActionContext struct {
	name    string
	sms     *model.StateMachines
	state   *model.SmState
	smTrans *model.SmTransition
}
type ChildSmExitActionContext struct {
	name    string
	sms     *model.StateMachines
	state   *model.SmState
	smTrans *model.SmTransition
}

func (smCtx *SmActionContext) setSmTransitionAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	//smCtx, ok := ruleCtx.(*SmActionContext)
	//if !ok {
	//	fmt.Printf("incorrect rule context type")
	//	return
	//}

	smTuple := tuples[model.TupleType(smCtx.name)]
	if smTuple == nil {
		fmt.Printf("sm not found %s", smCtx.name)
		return
	}

	smt, ok := smTuple.(model.StateMachineTuple)
	if !ok {
		fmt.Printf("%s not of type statemachinetuple", smCtx.name)
		return
	}

	fmt.Printf("setting sm[%s] to next state [%s] from state [%s]\n",
		smt.GetKey().String(), smCtx.smTrans.ToState, smt.GetState())
	smt.SetState(ctx, smCtx.smTrans.ToState)
	_ = startTimeoutForCurrentState(smt, rs)
}

func StartSm(ctx context.Context, rs model.RuleSession, s model.StateMachineTuple) error {
	if s.IsStarted() {
		return nil
	}
	_ = startTimeoutForCurrentState(s, rs)

	err := rs.Assert(ctx, s)
	if err != nil {
		return err
	}
	s.SetStarted(true)
	return nil
}

//to be called right after changing state to nextState
func startTimeoutForCurrentState(s model.StateMachineTuple, rs model.RuleSession) error {
	timer := s.GetStateTimeoutTimer()
	if timer != nil {
		timer.Stop()
	}
	smm := s.GetStateMachine()
	stt := smm.GetSmForState(s.GetState())
	if stt == nil {
		fmt.Printf("state transitions not found for state [%s]\n", s.GetState())
		return nil
	}
	timeout := stt.Timeout
	if timeout > 0 {
		timer = time.NewTimer(time.Duration(timeout) * time.Second)
		s.SetStateTimeoutTimer(timer)

		go func() {
			<-timer.C
			vals := s.GetMap()["sm_key"].(string)
			assertTimerTuple(rs, vals)
		}()
	}
	return nil
}

func assertTimerTuple(rs model.RuleSession, smKey string) {
	now := time.Now().UnixNano()

	timer, _ := model.NewTupleWithKeyValues("timer", now)
	_ = timer.SetString(context.TODO(), "ctx", smKey)

	_ = rs.Assert(context.TODO(), timer)

}
func (smCtx *ChildSmExitActionContext) childSmExitAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	//todo: change state to next state..
}

func (smCtx *StartChildSmActionContext) startChildSmAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	//smCtx, ok := ruleCtx.(*SmActionContext)
	//if !ok {
	//	fmt.Printf("incorrect rule context type")
	//	return
	//}

	smTuple := tuples[model.TupleType(smCtx.name)]
	if smTuple == nil {
		fmt.Printf("sm not found %s", smCtx.name)
		return
	}

	smt, ok := smTuple.(model.StateMachineTuple)
	if !ok {
		fmt.Printf("%s not of type statemachinetuple", smCtx.name)
		return
	}

	fmt.Printf("setting sm[%s] to next state [%s] from state [%s]\n",
		smt.GetKey().String(), smCtx.smTrans.ToState, smt.GetState())
	smt.SetState(ctx, smCtx.smTrans.ToState)
	_ = startTimeoutForCurrentState(smt, rs)
}
