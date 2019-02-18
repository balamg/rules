import pyrules

def MyConditionCb (ruleName, condiName, tupleJson):
    print ("In MyPython Condition\n");
    return 1;

def MyActionCb (ruleName, tupleJson):
    print ("In MyPython Action\n");
    return 0;


pyrules.RegisterTupleDescriptorsFromFile("/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json")

pyrules.CreateRuleSession ("rs")

pyrules.StartRuleSession("rs")

pyrules.AddRule("rs", "myrule1", "idr", MyConditionCb, MyActionCb)

pyrules.Assert("rs", "xyz")



