import pyrules

def MyConditionCbFromJson(ruleName, conditionName, tupleJsonStr):
    print ("In MyPython Condition\n");
    print(ruleName, conditionName, tupleJsonStr);
    tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)

    n1 = tupleMap["n1"]
    age = n1.Tuples["age"]

    n2 = tupleMap["n2"]
    age2 = n2.Tuples["age"]

    return age == age2

def MyConditionCb(ruleName, conditionName, tupleMap):
    print ("In MyPython Condition\n");
    n1 = tupleMap["n1"]
    age = n1.Tuples["age"]

    n2 = tupleMap["n2"]
    age2 = n2.Tuples["age"]

    return age == age2

def MyActionCbFromJson(ruleName, tupleJsonStr):
    print ("In MyPython Action\n");
    tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
    n1 = tupleMap["n1"]
    nm1 = n1.Tuples["name"]
    n2 = tupleMap["n2"]
    nm2 = n2.Tuples["name"]

    print nm1, nm2, tupleMap["n3"].Tuples["name"]

def MyActionCb(ruleName, tupleMap):
    print ("In MyPython Action\n");
    n1 = tupleMap["n1"]
    nm1 = n1.Tuples["name"]
    n2 = tupleMap["n2"]
    nm2 = n2.Tuples["name"]

    print nm1, nm2, tupleMap["n3"].Tuples["name"]

if __name__ == "__main__":

    pyrules.RegisterTupleDescriptorsFromFile("pyrulesapp.json")

    pyrules.CreateRuleSession ("rs")

    pyrules.StartRuleSession("rs")

    pyrules.AddRule("rs", "myrule1", ["n1", "n2", "n3"], MyConditionCb, MyActionCb)


    #Construct a tuple from a map
    props = {}
    props['name'] = "Bob"
    props['age'] = 48
    props['gender'] = "Male"
    props['salary'] = 100.1
    tuple = pyrules.Tuple("n1", props)

    pyrules.AssertTuple("rs", tuple)

    props = {}
    props['name'] = "Tom"
    props['age'] = 48
    props['gender'] = "Male"
    props['salary'] = 100.1
    tuple = pyrules.Tuple("n2", props)
    pyrules.AssertTuple("rs", tuple)

    props = {}
    props['name'] = "Pete"
    props['age'] = 48
    props['gender'] = "Male"
    props['salary'] = 100.1
    tuple = pyrules.Tuple("n3", props)
    pyrules.AssertTuple("rs", tuple)

