#ifndef _PYWRAP_H
#define _PYWRAP_H

int EvalRuleCondition (const char *ruleConditionModuleName, const char *ruleConditionFnName, const char *ruleName,
    const char *conditionName, const char *tupleJson);
int EvalRuleAction (const char *ruleActionFnModule, const char *ruleActionFnName, const char *ruleName, const char *tupleJson);
#endif
