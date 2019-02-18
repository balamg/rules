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
	fmt.Printf("Registered tuple descriptors.\n")
}

//export CreateRuleSession
func CreateRuleSession(ruleSessionName string) {
	_, found := ruleSessions[ruleSessionName]
	if !found {
		rs, _ := ruleapi.GetOrCreateRuleSession(ruleSessionName)
		ruleSessions[ruleSessionName] = rs
	}
	fmt.Printf("Created rulesession [%s].\n", ruleSessionName)
}

//export StartRuleSession
func StartRuleSession (ruleSessionName string) {
	rs := ruleSessions[ruleSessionName]
	rs.Start(nil)
	fmt.Printf("Started rulesession [%s].\n", ruleSessionName)
}

//export AddRule
func AddRule(ruleSessionName string, ruleName string, idrJson string) {
	//fmt.Printf("Adding a rule..[%s][%s][%s]\n",ruleSessionName, ruleName, idrJson)
	rule := ruleapi.NewRule(ruleName)
	rule.AddCondition("c1", []string{"n1"}, CCondition, nil)
	rule.SetAction(CAction)
	rs := ruleSessions[ruleSessionName]
	rs.AddRule(rule)
	fmt.Printf("Rule [%s] added to rulesession [%s].\n", ruleName, ruleSessionName)
}

//export Assert
func Assert(ruleSessionName string, tupleJson string) {
	//Now assert a "n1" tuple
	fmt.Printf("Asserting tuple [%s].\n", tupleJson)
	t2, _ := model.NewTupleWithKeyValues("n1", "Bob")
	t2.SetString(nil, "name", "Bob")
	rs := ruleSessions[ruleSessionName]
	rs.Assert(nil, t2)
}


func CAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tupleJson := ""
	//fmt.Printf("CAction started..[%s][%s]\n", ruleName, tupleJson)
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	//fmt.Printf("Context is [%s]\n", ruleCtx)
	C.performAction(C.CString(ruleName), C.CString(tupleJson))
	//fmt.Printf("CAction complete..\n")
}

func CCondition(condName string, ruleName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//tupleJson := "x"
	//fmt.Printf("CCondition started..[%s][%s][%s]\n", ruleName, condName, "")
	//var ret bool
	i := C.evalCondition(C.CString(ruleName), C.CString(condName), C.CString(""))
	//fmt.Printf("CCondition complete..[%d]\n", i)
	return i != 0
}

func main() {}
