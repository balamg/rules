if __name__ == "__main__":
    import pyrules

    def MyConditionCb (ruleName, conditionName, tupleMap):
        print ("In MyPython Condition\n");
        tuple = tupleMap["n1"]
        age = tuple.Tuples["age"]
        print age
        return 1

    def MyActionCb (ruleName, tupleMap):
        print ("In MyPython Action\n");
        tuple = tupleMap["n1"]
        gender = tuple.Tuples["gender"]
        print gender
        return 1

    pyrules.RegisterTupleDescriptorsFromFile("/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json")

    pyrules.CreateRuleSession ("rs")

    pyrules.StartRuleSession("rs")

    pyrules.AddRule("rs", "myrule1", "idr", MyConditionCb, MyActionCb)


    #Construct a tuple from a map
    props = {}
    props['name'] = "Bala"
    props['age'] = 48
    props['gender'] = "Male"
    props['salary'] = 100.1
    tuple = pyrules.Tuple("n1", props)

    pyrules.AssertTuple("rs", tuple)



