package main
// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "actionconditioncb.h"
import "C"

import (
	"github.com/project-flogo/rules/common/model"
	"fmt"
	"context"
)


//export RegisterTupleDescriptors
func RegisterTupleDescriptors(tupleDescriptor string) {
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}
}

func CAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("PyAction..\n")
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	fmt.Printf("Context is [%s]\n", ruleCtx)
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return
	}
	C.callPyCb(C.CString(ruleName))
	fmt.Printf("PyAction..done\n")
}

func CCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	fmt.Printf("PyCondition..\n")
	//This conditions filters on name="Bob"
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name, _ := t1.GetString("name")
	C.callPyCb(C.CString(condName))
	fmt.Printf("PyCondition..done\n")
	return name == "Bob"
}

func main() {}
