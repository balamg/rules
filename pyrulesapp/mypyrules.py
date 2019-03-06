import pyrules

def c_nametom(ruleName, conditionName, tupleJsonStr):
    try:
        print ("In c_nametom\n");
        print(ruleName, conditionName, tupleJsonStr);

        tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
        n1 = tupleMap["n1"]
        name = n1.Tuples["name"]
        return name == "Tom"
    except (Exception):
        return 0


def nametom(ruleName, tupleJsonStr):
    try:
        print ("In nametom\n");
        tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
        n1 = tupleMap["n1"]
        nm1 = n1.Tuples["name"]
        print nm1
    except (Exception):
        return

def c_bothnamestom(ruleName, conditionName, tupleJsonStr):
    try:
        print ("In c_bothnamestom\n");
        print(ruleName, conditionName, tupleJsonStr);

        tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
        n1 = tupleMap["n1"]
        name1 = n1.Tuples["name"]

        n2 = tupleMap["n2"]
        name2 = n2.Tuples["name"]

        print name1, name2
        return (name1 == "Tom" and name1 == name2)

    except (Exception):
        return 0


def bothnamestom(ruleName, tupleJsonStr):
    try:
        print ("In bothnamestom\n");
        tupleMap = pyrules.TuplesFromJsonStr(tupleJsonStr)
        n1 = tupleMap["n1"]
        nm1 = n1.Tuples["name"]
        print nm1
    except (Exception):
        return