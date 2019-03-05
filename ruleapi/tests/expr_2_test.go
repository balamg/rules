package tests

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"golang.org/x/net/context"
	"testing"
)

//2 conditions, 1 expr each
func Test_2_Expr(t *testing.T) {

	rs, _ := createRuleSession()
	r1 := ruleapi.NewRule("r1")
	r1.AddExprCondition("c1", "$.t1.p1 > $.t2.p1", nil)
	r1.AddExprCondition("c2", "$.t1.p1 == 2", nil)
	r1.SetAction(a2)
	rs.AddRule(r1)

	rs.Start(nil)

	var ctx context.Context

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	t1.SetInt(nil, "p1", 2)
	t1.SetDouble(nil, "p2", 1.3)
	t1.SetString(nil, "p3", "t3")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t1)

	t2, _ := model.NewTupleWithKeyValues("t2", "t2")
	t2.SetInt(nil, "p1", 1)
	t2.SetDouble(nil, "p2", 1.1)
	t2.SetString(nil, "p3", "t3")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t2)
	rs.Unregister()
}

func a2(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Test_2_Expr executed!")
}