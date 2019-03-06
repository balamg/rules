from ctypes import *
import pyrules

conditionDict = {}
actionDict = {}


class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

def GoStr (str):
    return GoString(str, len(str))


def PyConditionCb (ruleName, conditionName, tupleJsonStr):
    key = "rs:" + ruleName + ":" + "c1"
    # print ("In PyConditionCb", tupleJsonStr);
    tuples = pyrules.TuplesFromJsonStr(tupleJsonStr)

    conditionCb = conditionDict[key]
    ret = conditionCb(ruleName, conditionName, tuples)
    return ret;

def PyActionCb (ruleName, tupleJsonStr):
    key = "rs:" + ruleName
    # print ("In PyActionCb", tupleJsonStr);
    tuples = pyrules.TuplesFromJsonStr(tupleJsonStr)

    actionCb = actionDict[key]
    actionCb(ruleName, tuples)
    return 0;


