package tests

import (
	"context"
	"fmt"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"testing"
)

func Test_ThreeConditions (t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("3-tuples")
	rule.AddCondition("c1", []string{"t3", "t4", "t5"}, three_tuples, nil)
	rule.SetAction(three_tuple_action)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())
	rs.Start(nil)

	ctx := context.WithValue(context.TODO(), TestKey{}, t)
	t1, _ := model.NewTupleWithKeyValues("t3", "t3")
	rs.Assert(ctx, t1)

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	t2, _ := model.NewTupleWithKeyValues("t4", "t4")
	rs.Assert(ctx, t2)

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	t3, _ := model.NewTupleWithKeyValues("t5", "t5")
	rs.Assert(ctx, t3)

	rs.Unregister()

}

func three_tuples(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	fmt.Printf("Condition fired: [%s], [%d]\n", ruleName,len(tuples))
	return true
}

func three_tuple_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Rule fired: [%s], [%d]\n", ruleName,len(tuples))
}
