import capi
import json
import pyrules

def RegisterTupleDescriptorsFromFile (tupleDescFileName):
    tupleDescJson = open(tupleDescFileName, 'r').read()
    RegisterTupleDescriptors (tupleDescJson)

def RegisterTupleDescriptors (tupleDescJson):
    capi.lib.RegisterTupleDescriptors(capi.GoStr(tupleDescJson))

def CreateRuleSession (ruleSessionName):
    capi.lib.CreateRuleSession(capi.GoStr(ruleSessionName))

def AddRule (ruleSessionName, ruleName, tupleTypesArr, condFn, actionFn):

    key = ruleSessionName + ":" + ruleName + ":" + "c1"
    capi.conditionDict[key] = condFn

    key = ruleSessionName + ":" + ruleName
    capi.actionDict[key] = actionFn

    tupleTypesJsonStr = json.dumps(tupleTypesArr)
    capi.lib.AddRule (capi.GoStr(ruleSessionName), capi.GoStr(ruleName), capi.GoStr(tupleTypesJsonStr))

def StartRuleSession (ruleSessionName):
    capi.lib.StartRuleSession(capi.GoStr(ruleSessionName))

def AssertTuple (ruleSessionName, tuple):
    tupleJson = pyrules.TuplesToJsonStr(tuple)
    capi.lib.Assert(capi.GoStr(ruleSessionName), capi.GoStr(tupleJson))

def Assert (ruleSessionName, tupleJson):
    capi.lib.Assert(capi.GoStr(ruleSessionName), capi.GoStr(tupleJson))


