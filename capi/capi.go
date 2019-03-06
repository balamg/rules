package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "actionconditioncb.h"
import "C"

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/project-flogo/rules/common/model"
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
func StartRuleSession(ruleSessionName string) {
	rs := ruleSessions[ruleSessionName]
	rs.Start(nil)
	fmt.Printf("Started rulesession [%s].\n", ruleSessionName)
}

//export AddRule
func AddRule(ruleSessionName string, ruleName string, tupleTypesJsonStr string) {
	fmt.Printf("Adding a rule..[%s][%s][%s]\n", ruleSessionName, ruleName, tupleTypesJsonStr)
	rule := ruleapi.NewRule(ruleName)
	tupleTypes := []string{}
	json.Unmarshal([]byte(tupleTypesJsonStr), &tupleTypes)

	rule.AddCondition("c1", tupleTypes, CCondition, nil)
	rule.SetAction(CAction)
	rs := ruleSessions[ruleSessionName]
	rs.AddRule(rule)
	fmt.Printf("Rule [%s] added to rulesession [%s].\n", ruleName, ruleSessionName)
}

//export Assert
func Assert(ruleSessionName string, tupleJson string) {
	//fmt.Printf("Got json: %s\n", tupleJson)
	//Now assert a "n1" tuple
	tuple := model.TupleFromJsonStr(tupleJson)
	rs := ruleSessions[ruleSessionName]
	fmt.Printf("Asserting tuple [%s].\n", tupleJson)
	rs.Assert(nil, tuple)
}

func CAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tuplesJson, err := json.Marshal(tuples)
	if err != nil {
		fmt.Printf("Error in CAction serialize: %s", err)
		return
	}
	//fmt.Println(string(tuplesJson))
	C.performAction(C.CString(ruleName), C.CString(string(tuplesJson)))
}

func CCondition(condName string, ruleName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	tuplesJson, err := json.Marshal(tuples)
	if err != nil {
		fmt.Printf("Error in CCondition serialize: %s", err)
		return false
	}
	//fmt.Printf("In CCondition: [%s]\n", string(tuplesJson))
	i := C.evalCondition(C.CString(ruleName), C.CString(condName), C.CString(string(tuplesJson)))
	//fmt.Printf("CCondition complete..[%d]\n", i)
	return i != 0
}

func main() {}
