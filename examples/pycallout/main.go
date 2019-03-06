package main

import "github.com/project-flogo/rules/pyembed"

/**
To run this program, ensure you set PYTHONPATH=$GOPATH/github.com/project-flogo/rules/examples/pyembed:
$GOPATH/github.com/project-flogo/rules
 */

func main() {
	jsonStr := "{\"n1\":{\"TupleType\":\"n1\",\"Tuples\":{\"age\":48,\"gender\":\"Male\",\"name\":\"Bob\",\"salary\":100.1}},\"n2\":{\"TupleType\":\"n2\",\"Tuples\":{\"age\":48,\"gender\":\"Male\",\"name\":\"Tom\",\"salary\":100.1}},\"n3\":{\"TupleType\":\"n3\",\"Tuples\":{\"age\":48,\"gender\":\"Male\",\"name\":\"Pete\",\"salary\":100.1}}}"
	isTrue := pyembed.EvalRuleCondition("pyrulesapp", "MyConditionCbFromJson", "rule1", "condition1", jsonStr)

	if isTrue {
		pyembed.EvalRuleAction("pyrulesapp", "MyActionCbFromJson", "rule1", jsonStr)
	}
}
