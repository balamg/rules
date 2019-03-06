package main

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/pyembed"
	"github.com/project-flogo/rules/ruleapi"
)

//
//TODO: To run this program, set PYTHONPATH to rules, rules/pyrules, rules/pyruleapp
//

func CAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tuplesJson, err := json.Marshal(tuples)
	if err != nil {
		fmt.Printf("Error in CAction serialize: %s", err)
		return
	}
	fmt.Println("action: " + string(tuplesJson))
	//C.performAction(C.CString(ruleName), C.CString(string(tuplesJson)))
	if ruleName == "nametom" {
		pyembed.EvalRuleAction("mypyrules", "nametom", ruleName, string(tuplesJson))
	} else if ruleName == "bothnamestom" {
		pyembed.EvalRuleAction("mypyrules", "bothnamestom", ruleName, string(tuplesJson))
	}
}

func CCondition(condName string, ruleName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	tf := false
	tuplesJson, _ := json.Marshal(tuples)

	fmt.Println("here.." + string(tuplesJson))
	if ruleName == "nametom" {
		tf = pyembed.EvalRuleCondition("mypyrules", "c_nametom", ruleName, condName, string(tuplesJson))
	} else if ruleName == "bothnamestom" {
		tf = pyembed.EvalRuleCondition("mypyrules", "c_bothnamestom", ruleName, condName, string(tuplesJson))

	}
	return tf
}

func init() {

	config.RegisterConditionEvaluator("c.nametom", CCondition)
	config.RegisterActionFunction("a.nametom", CAction)

	config.RegisterConditionEvaluator("c.bothnamestom", CCondition)
	config.RegisterActionFunction("a.bothnamestom", CAction)

}

func main() {

	fmt.Println("** rulesapp: Example usage of the Rules module/API **")

	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetAbsPathForResource("src/github.com/project-flogo/rules/pyrulesapp/pyrulesapp_types.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	//Create a RuleSession
	//rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	ruleConfigFile := common.GetAbsPathForResource("src/github.com/project-flogo/rules/pyrulesapp/pyrules.json")
	ruleConfigJson := common.FileToString(ruleConfigFile)

	rs, err := ruleapi.GetOrCreateRuleSessionFromConfig("rs", string(ruleConfigJson))

	if err != nil {
		fmt.Printf("Error [%s]\n", err)
	}

	rs.Start(nil)

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Tom")
	t1, err := model.NewTupleWithKeyValues("n1", "Tom")
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
	}
	t1.SetString(nil, "name", "Tom")
	err = rs.Assert(nil, t1)

	if err != nil {
		fmt.Printf("Error [%s]\n", err)
	}

	//Now assert a "n1" tuple
	fmt.Println("Asserting n2 tuple with name=Tom")
	t2, err := model.NewTupleWithKeyValues("n2", "Tom")
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
	}
	t1.SetString(nil, "name", "Tom")
	err = rs.Assert(nil, t2)

	fmt.Printf("Done..\n")
}
