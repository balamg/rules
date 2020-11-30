package ruleapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"

	"github.com/ghodss/yaml"
	"github.com/project-flogo/rules/common"
)

func Test_SM(t *testing.T) {

	var sm []model.StateMachine
	smStr := common.FileToString("sm.yaml")
	if smStr == "" {
		t.FailNow()
	}
	err := yaml.Unmarshal([]byte(smStr), &sm)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}
	fmt.Printf("%v\n", sm)

	//CreateRule(sm)

	err = model.RegisterSmTypes(sm)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	rs, err := GetOrCreateRuleSession("asession")
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	for i := range sm {
		rules, err := CreateRulesForSm(sm[i])
		if err != nil {
			fmt.Printf("%s", err)
			t.FailNow()
		}

		for i := range rules {
			err := rs.AddRule(rules[i])
			if err != nil {
				fmt.Printf("%s", err)
				t.FailNow()
			}
		}
	}
	fmt.Printf("%v\n", sm)

	valMap := map[string]interface{}{"sm_key": "s1"}

	smt, err := model.NewStateMachineTuple(sm[0], valMap)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}
	err = rs.Start(nil)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	//err = rs.Assert(context.TODO(), smt)
	//if err != nil {
	//	fmt.Printf("%s", err)
	//	t.FailNow()
	//}
	err = StartSm(context.TODO(), rs, smt)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	time.Sleep(20 * time.Second)
}
