package model

import (
	"context"
	"fmt"
	"testing"

	"github.com/project-flogo/rules/ruleapi"
)

func Test_SM(t *testing.T) {

	sm := &SmState{
		Transitions: map[string]*SmTransition{"c1": &SmTransition{
			"exitS1", "S2", "enterS1",
		}},
		Timeout:      10,
		TimeoutState: "s1t0",
	}

	CreateRule(sm)

}

func CreateRule(sm *SmState) {

	for condition, smTrans := range sm.Transitions {
		ruleName := fmt.Sprintf("%s_%s_%s", sm.Name, condition, smTrans.NextState)

		rule := ruleapi.NewRule(ruleName)
		sa := SmActionContext{condition: condition, smTrans: smTrans}
		rule.SetAction(sa.setSmTransitionAction)
		rule.AddExprCondition(condition, condition, nil)
		rule.SetPriority(1)
	}

	ruleName := fmt.Sprintf("%s_timeout", sm.Name)

	rule := ruleapi.NewRule(ruleName)
	rule.AddExprCondition("timeout", "$.timer['v1'] == ..", nil)
	smt := &SmTimeoutActionContext{sm}
	rule.SetAction(smt.TimeoutAction)
}

func (smt *SmTimeoutActionContext) TimeoutAction(ctx context.Context, session RuleSession, s string, m map[TupleType]Tuple, ruleContext RuleContext) {
	//todo:
	//get the sm tuple
	//and set its state to next state
	//smt.sm.TimeoutState
}

type SmTimeoutActionContext struct {
	sm *SmState
}
type SmActionContext struct {
	condition string
	smTrans   *SmTransition
}

func (sm *SmActionContext) setSmTransitionAction(ctx context.Context, session RuleSession, ruleName string, tuples map[TupleType]Tuple, ruleCtx RuleContext) {

}
