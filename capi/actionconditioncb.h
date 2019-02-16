#ifndef _ACTIONCONDITIONCB_H
#define _ACTIONCONDITIONCB_H

int registerConditionCb (int (*conditionCbf)(const char *ruleName, const char *conditionName, const char* tupleJson));
int registerActionCb (int (*actionCbf)(const char *ruleName, const char* tupleJson));

int evalCondition(const char *ruleName, const char *conditionName, const char *tupleJson);
int performAction(const char *ruleName, const char *tupleJson);

#endif
