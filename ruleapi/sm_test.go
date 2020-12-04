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

	var sms model.StateMachines
	smStr := common.FileToString("sm.yaml")
	if smStr == "" {
		t.FailNow()
	}
	err := yaml.Unmarshal([]byte(smStr), &sms)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
	//fmt.Printf("%v\n", sms)

	//CreateRule(sm)

	err = model.RegisterSmTypes(sms.StateMachines)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	rs, err := GetOrCreateRuleSession("asession")
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	for i := range sms.StateMachines {
		rules, err := CreateRulesForSm(&sms, &sms.StateMachines[i])
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
			fmt.Printf("added rule: [%s]\n", rules[i].GetName())
		}
	}
	fmt.Printf("Added rules successfully..\n")

	valMap := map[string]interface{}{"sm_key": "s1"}

	smt, err := model.NewStateMachineTuple(&sms.StateMachines[0], valMap)
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
	err = StartSm(context.TODO(), rs, &sms.StateMachines[0], smt)
	if err != nil {
		fmt.Printf("%s", err)
		t.FailNow()
	}

	time.Sleep(20 * time.Second)
}
