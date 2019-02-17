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
	fmt.Printf("Registered tuple descriptors\n")
}

//export CreateRuleSession
func CreateRuleSession(name string) {
	_, found := ruleSessions[name]
	if !found {
		rs, _ := ruleapi.GetOrCreateRuleSession(name)
		ruleSessions[name] = rs
	}
}

//export StartRuleSession
func StartRuleSession (ruleSessionName string) {
	rs := ruleSessions[ruleSessionName]
	rs.Start(nil)
}

//export AddRule
func AddRule(ruleSessionName string, ruleName string, idrJson string) {
	//fmt.Printf("Adding a rule..[%s][%s][%s]\n",ruleSessionName, ruleName, idrJson)
	rule := ruleapi.NewRule("n1.name == Bob")
	rule.AddCondition("c1", []string{"n1"}, CCondition, nil)
	rule.SetAction(CAction)
	rule.SetContext("This is a test of context")
	rs := ruleSessions[ruleSessionName]
	rs.AddRule(rule)
	fmt.Printf("Adding a rule successful.\n")

}

//export Assert
func Assert(ruleSessionName string, tupleJson string) {
	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Bob\n")
	t2, _ := model.NewTupleWithKeyValues("n1", "Bob")
	t2.SetString(nil, "name", "Bob")
	rs := ruleSessions[ruleSessionName]
	rs.Assert(nil, t2)
	//fmt.Println("Asserted n1 tuple with name=Bob\n")
}


func CAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tupleJson := ""
	//fmt.Printf("CAction started..[%s][%s]\n", ruleName, tupleJson)
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	//fmt.Printf("Context is [%s]\n", ruleCtx)
	C.performAction(C.CString(ruleName), C.CString(tupleJson))
	//fmt.Printf("CAction complete..\n")
}

func CCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//tupleJson := "x"
	//fmt.Printf("CCondition started..[%s][%s][%s]\n", ruleName, condName, tupleJson)
	//var ret bool
	i := C.evalCondition(C.CString(ruleName), C.CString(condName), C.CString(""))
	//fmt.Printf("CCondition complete..\n")
	return i != 0
}

func main() {}
