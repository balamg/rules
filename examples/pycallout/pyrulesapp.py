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


def MyActionCbFromJson(ruleName, tupleJsonStr):
    print ("In MyPython Action\n");
    tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
    n1 = tupleMap["n1"]
    nm1 = n1.Tuples["name"]
    n2 = tupleMap["n2"]
    nm2 = n2.Tuples["name"]

    print nm1, nm2, tupleMap["n3"].Tuples["name"]