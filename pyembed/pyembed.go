package pyembed

//#cgo CFLAGS: -g -Wall
//#cgo pkg-config: python-2.7
// #include <stdlib.h>
// #include "pyembed.h"
// #include <Python.h>
import "C"

func EvalRuleCondition(ruleConditionModuleName string, ruleConditionFnName string, ruleName string,
	conditionName string, tuplesJson string) bool {
	i := C.EvalRuleCondition(C.CString(ruleConditionModuleName), C.CString(ruleConditionFnName),
			C.CString(ruleName), C.CString(conditionName), C.CString(tuplesJson))
	j := int(i)
	return j == 1 //all other values are treated as false
}

func EvalRuleAction(ruleActionModuleName string, ruleActionFnName string, ruleName string, tuplesJson string) int {
	i := C.EvalRuleAction(C.CString(ruleActionModuleName),C.CString(ruleActionFnName), C.CString(ruleName), C.CString(tuplesJson))
	j := int(i)
	return j
}
