package utils

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/project-flogo/core/data/property"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/rules/common/model"
)

var td tuplePropertyResolver
var resolver resolve.CompositeResolver
var factory expression.Factory
var re *regexp.Regexp
var parsedExpressions map[string]expression.Expr
var parsedExprLock sync.Mutex

func init() {
	td = tuplePropertyResolver{}
	resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{
		".":        &td,
		"env":      &resolve.EnvResolver{},
		"property": &property.Resolver{},
		"loop":     &resolve.LoopResolver{},
	})

	factory = script.NewExprFactory(resolver)
	re = regexp.MustCompile(`\$(\w+)`)
	parsedExpressions = map[string]expression.Expr{}

}

func Eval(cExpr string, tuple model.Tuple, toReplace bool) (interface{}, error) {
	//fmt.Printf("evaluating condition: [%s], [%s], [%s]\n", cExpr, condName, ruleNm)
	var result interface{}
	if cExpr != "" {
		if toReplace {
			cExpr = getModifiedExpr(cExpr)
		}
		var exprn expression.Expr
		var err error
		parsedExprLock.Lock()
		exprn = parsedExpressions[cExpr]
		if exprn == nil {
			exprn, err = factory.NewExpr(cExpr)
			if err != nil {
				fmt.Printf("Error forming expression: [%s]\n", cExpr)
				return 0, err
			}
			parsedExpressions[cExpr] = exprn
		}
		parsedExprLock.Unlock()
		tuples := map[model.TupleType]model.Tuple{}
		if toReplace {
			tuples["X"] = tuple
		} else {
			tuples[tuple.GetTupleType()] = tuple
		}
		scope := tupleScope{tuples}
		result, err = exprn.Eval(&scope)
		if err != nil {
			fmt.Printf("Error evaluating condition: [%s]\n", cExpr)
			return 0, err
		}
	}
	return result, nil
}

//////////////////////////////////////////////////////////
type tupleScope struct {
	tuples map[model.TupleType]model.Tuple
}

func (ts *tupleScope) GetValue(name string) (value interface{}, exists bool) {
	return false, true
}

func (ts *tupleScope) SetValue(name string, value interface{}) error {
	return nil
}

// SetAttrValue sets the value of the specified attribute
func (ts *tupleScope) SetAttrValue(name string, value interface{}) error {
	return nil
}

///////////////////////////////////////////////////////////
type tuplePropertyResolver struct {
}

func (t *tuplePropertyResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	ts := scope.(*tupleScope)
	tuple := ts.tuples[model.TupleType(field)]
	if tuple == nil {
		return nil, fmt.Errorf("Tuple [%s] not found in scope", field)
	} else {
		m := tuple.GetMap()
		return m, nil
	}
}

func (*tuplePropertyResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolve.NewResolverInfo(false, false)
}

func getModifiedExpr(expr string) string {
	return re.ReplaceAllString(expr, "$.X.$1")
}

func GetRefs(cstr string) []string {
	keys2 := []string{}
	keys := re.FindAllStringSubmatch(cstr, -1)
	for _, k := range keys {
		keys2 = append(keys2, k[1])
	}
	return keys2
}
