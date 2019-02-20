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

def RegisterTupleDescriptorsFromFile (tupleDescFileName):
    tupleDescJson = open(tupleDescFileName, 'r').read()
    RegisterTupleDescriptors (tupleDescJson)

def RegisterTupleDescriptors (tupleDescJson):
    pyrules.lib.RegisterTupleDescriptors(GoStr(tupleDescJson))

def CreateRuleSession (ruleSessionName):
    pyrules.lib.CreateRuleSession(GoStr(ruleSessionName))

def AddRule (ruleSessionName, ruleName, idrs, condFn, actionFn):

    key = ruleSessionName + ":" + ruleName + ":" + "c1"
    conditionDict[key] = condFn

    key = ruleSessionName + ":" + ruleName
    actionDict[key] = actionFn

    pyrules.lib.AddRule (GoStr(ruleSessionName), GoStr(ruleName), GoStr(idrs))

def StartRuleSession (ruleSessionName):
    pyrules.lib.StartRuleSession(GoStr(ruleSessionName))

def AssertTuple (ruleSessionName, tuple):
    tupleJson = pyrules.TuplesToJsonStr(tuple)
    pyrules.lib.Assert(GoStr(ruleSessionName),GoStr(tupleJson))

def Assert (ruleSessionName, tupleJson):
    pyrules.lib.Assert(GoStr(ruleSessionName),GoStr(tupleJson))


