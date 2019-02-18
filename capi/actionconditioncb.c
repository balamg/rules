#include "actionconditioncb.h"
#include <stdio.h>

//condition function pointer
int (*conditionCb)(const char *ruleName, const char *conditionName, const char* tupleJson);

//action function pointer
int (*actionCb)(const char *ruleName, const char* tupleJson);


int registerConditionCb (int (*conditionCbf)(const char *ruleName, const char *conditionName, const char* tupleJson)) {
    conditionCb = conditionCbf;
//    printf ("registered condition callback function....\n");
    return 0;
}

int registerActionCb (int (*actionCbf)(const char *ruleName, const char* tupleJson)) {
    actionCb = actionCbf;
//    printf ("registered action callback function....\n");
    return 0;
}

int evalCondition(const char *ruleName, const char *conditionName, const char *tupleJson) {
    if (conditionCb != NULL) {
        int x = conditionCb (ruleName, conditionName, tupleJson);
//         printf ("return from Py [%s][%s][%s]\n", ruleName, conditionName, tupleJson);
         return x;
    } else {
        printf ("condition function not registered\n");
    }
    return 0;
}


int performAction(const char *ruleName, const char *tupleJson) {
    if (actionCb != NULL) {
        return actionCb (ruleName, tupleJson);
    } else {
        printf ("action function not registered\n");
    }
    return 0;
}

