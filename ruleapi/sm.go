package ruleapi

import (
	"context"
	"fmt"
	"time"

	"github.com/project-flogo/rules/common/model"
)

func CreateRulesForSm(sm model.StateMachineModel) ([]model.Rule, error) {
	var rules []model.Rule

	for i := range sm.States {
		rls, err := CreateRulesForState(sm.Name, sm.States[i])
		if err != nil {
			return nil, err
		}
		for i := range rls {
			rules = append(rules, rls[i])
		}
	}

	return rules, nil
}

func CreateRulesForState(smName string, sm model.SmState) ([]model.Rule, error) {

	var rules []model.Rule

	for i := range sm.Transitions {
		smTrans := sm.Transitions[i]
		condition := sm.Transitions[i].Condition
		ruleName := fmt.Sprintf("%s_%s_%s", sm.State, condition, smTrans.ToState)

		rule := NewRule(ruleName)
		sa := SmActionContext{name: smName, smTrans: &smTrans}
		rule.SetContext(&sa)
		rule.SetAction(sa.setSmTransitionAction)

		currStateCondition := fmt.Sprintf("$.%s.sm_state == '%s'", smName, sm.State)
		err := rule.AddExprCondition(currStateCondition, currStateCondition, &smTrans)
		if err != nil {
			return rules, err
		}
		err = rule.AddExprCondition(condition, condition, nil)
		if err != nil {
			return rules, err
		}
		rule.SetPriority(1)
		rules = append(rules, rule)
	}
	if sm.Timeout > 0 {
		timeoutRule, err := setTimeoutRuleForState(smName, sm)
		if err != nil {
			return rules, err
		}
		rules = append(rules, timeoutRule)
	}
	return rules, nil
}

func setTimeoutRuleForState(smName string, state model.SmState) (model.Rule, error) {
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

	smt := &SmTimeoutActionContext{smName, &state}
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
