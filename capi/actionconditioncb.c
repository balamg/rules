#include "actionconditioncb.h"
#include <stdio.h>

//condition function pointer
int (*conditionCb)(const char *ruleName, const char *conditionName, const char* tupleJson);

//action function pointer
int (*actionCb)(const char *ruleName, const char* tupleJson);

int my_fun (const char *name) {
    printf ("Hi there %s\n", name);
    return 0;
}
int registerConditionCb (int (*conditionCbf)(const char *ruleName, const char *conditionName, const char* tupleJson)) {
    conditionCb = conditionCbf;
    printf ("registered condition callback function....\n");
    return 0;
}

int registerActionCb (int (*actionCbf)(const char *ruleName, const char* tupleJson)) {
    actionCb = actionCbf;
    printf ("registered action callback function....\n");
    return 0;
}

int evalCondition(const char *ruleName, const char *conditionName, const char *tupleJson) {
    if (conditionCb != NULL) {
        conditionCb (ruleName, conditionName, tupleJson);
    } else {
        printf ("condition function not registered\n");
    }
    return 0;
}


int performAction(const char *ruleName, const char *tupleJson) {
    if (actionCb != NULL) {
        actionCb (ruleName, tupleJson);
    } else {
        printf ("action function not registered\n");
    }
    return 0;
}

