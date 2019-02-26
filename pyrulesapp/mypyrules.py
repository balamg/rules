import pyrules

def MyConditionCbFromJson(ruleName, conditionName, tupleJsonStr):
    try:
        print ("In MyConditionCbFromJson\n");
        print(ruleName, conditionName, tupleJsonStr);

        tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
        n1 = tupleMap["n1"]
        name = n1.Tuples["name"]
        return name == "Tom"
    except (Exception):
        return 0


def MyActionCbFromJson(ruleName, tupleJsonStr):
    try:
        print ("In MyActionCbFromJson\n");
        tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
        n1 = tupleMap["n1"]
        nm1 = n1.Tuples["name"]
        print nm1
    except (Exception):
        return