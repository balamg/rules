package main
// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "actionconditioncb.h"
import "C"

import (
	"github.com/project-flogo/rules/common/model"
	"fmt"
	"context"
	"github.com/project-flogo/rules/ruleapi"
)
var (
	ruleSessions map[string]model.RuleSession
)

func init() {
	ruleSessions = make(map[string]model.RuleSession)
}
//export RegisterTupleDescriptors
func RegisterTupleDescriptors(tupleDescriptor string) {
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}
}

//export CreateRuleSession
func CreateRuleSession(name string) {
	_, found := ruleSessions[name]
	if !found {
		rs, _ := ruleapi.GetOrCreateRuleSession(name)
		ruleSessions[name] = rs
	}
}

//export AddRule
func AddRule(ruleSessionName string, ruleName string, idr[]string) {
	rule := ruleapi.NewRule("n1.name == Bob")
	rule.AddCondition("c1", []string{"n1"}, CCondition, nil)
	rule.SetAction(CAction)
	rule.SetContext("This is a test of context")
	rs := ruleSessions[ruleSessionName]
	rs.AddRule(rule)
}

//export Assert
func Assert(ruleSessionName string, tupleJson string) (err error) {

	return nil
}


func CAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("CAction started..\n")
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	fmt.Printf("Context is [%s]\n", ruleCtx)
	tupleJson := ""
	C.performAction(C.CString(ruleName), C.CString(tupleJson))
	fmt.Printf("CAction complete..\n")
}

func CCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	fmt.Printf("CCondition started..\n")
	var ret bool
	tupleJson := ""
	C.evalCondition(C.CString(ruleName), C.CString(condName), C.CString(tupleJson))
	fmt.Printf("CCondition complete..\n")
	return ret
}

func main() {}
