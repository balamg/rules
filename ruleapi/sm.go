package ruleapi

import (
	"context"
	"fmt"

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
		sa := SmActionContext{name: smName, condition: condition, smTrans: &smTrans}
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

	//ruleName := fmt.Sprintf("%s_timeout", sm.State)
	//
	//rule := NewRule(ruleName)
	//rule.AddExprCondition("timeout", "$.timer['v1'] == ..", nil)
	//smt := &SmTimeoutActionContext{&sm}
	//rule.SetAction(smt.TimeoutAction)
	return rules, nil
}

func (smt *SmTimeoutActionContext) TimeoutAction(ctx context.Context, session model.RuleSession, s string, m map[model.TupleType]model.Tuple, ruleContext model.RuleContext) {
	//todo:
	//get the sm tuple
	//and set its state to next state
	//smt.sm.TimeoutState
}

type SmTimeoutActionContext struct {
	sm *model.SmState
}
type SmActionContext struct {
	name      string
	condition string
	smTrans   *model.SmTransition
}

func (sm *SmActionContext) setSmTransitionAction(ctx context.Context, session model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	smCtx, ok := ruleCtx.(*SmActionContext)
	if !ok {
		fmt.Printf("incorrect rule context type")
		return
	}

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
