#include <Python.h>
#include "pyembed.h"

#include <stdio.h>

int EvalRuleCondition (const char *ruleConditionModuleName, const char *ruleConditionFnName, const char *ruleName, const char *conditionName, const char *tupleJson) {
      const char* p = Py_GetPath();
   printf("Initializing python %s\n", p);
   Py_Initialize();
   if (!Py_IsInitialized()) {
        printf("Unable to initialize Python interpreter.");
        Py_Finalize();
        return -1;
   }
   PyObject *pyMod = PyImport_ImportModule(ruleConditionModuleName);
   if (pyMod == NULL) {
        printf ("Error importing [%s]\n", ruleConditionModuleName);
        Py_Finalize();
        return -1;
   }

   PyObject *pyFunc = NULL;
   pyFunc = PyObject_GetAttrString(pyMod, ruleConditionFnName);
   if (pyFunc == NULL) {
        printf ("Error importing [%s]\n", ruleConditionFnName);
        Py_DECREF(pyMod);
        Py_DECREF(pyFunc);
        Py_Finalize();
        return -1;
   }
   Py_DECREF(pyMod);

   PyObject *retVal;
   int iRetVal = 0;
   retVal = PyEval_CallFunction(pyFunc, "sss", ruleName, conditionName, tupleJson);
   if (retVal == NULL) {
        printf ("Error invoking function [%s]\n", ruleConditionFnName);
        Py_DECREF(pyFunc);
        Py_DECREF(retVal);
        Py_Finalize();
        return -1;
   }
   Py_DECREF(pyFunc);

   PyArg_Parse(retVal, "i", &iRetVal);
   Py_DECREF(retVal);

   printf ("Result [%d]\n", iRetVal);

   Py_Finalize();

   return iRetVal;
}

int EvalRuleAction (const char *ruleActionFnModule, const char *ruleActionFnName, const char *ruleName, const char *tupleJson) {
   printf("Initializing python\n");
   Py_Initialize();
   if (!Py_IsInitialized()) {
        printf("Unable to initialize Python interpreter.");
        Py_Finalize();
        return -1;
   }
   PyObject *pyMod = PyImport_ImportModule(ruleActionFnModule);
   if (pyMod == NULL) {
        printf ("Error importing [%s]\n", ruleActionFnModule);
        Py_DECREF(pyMod);
        Py_Finalize();
        return -1;
   }


   PyObject *pyFunc = NULL;
   pyFunc = PyObject_GetAttrString(pyMod, ruleActionFnName);
   if (pyFunc == NULL) {
        printf ("Error importing [%s]\n", ruleActionFnName);
        Py_DECREF(pyMod);
        Py_DECREF(pyFunc);
        Py_Finalize();
        return -1;
   }
   Py_DECREF(pyMod);


   PyObject *retVal;
   int iRetVal;
   retVal = PyEval_CallFunction(pyFunc, "ss", ruleName, tupleJson);
   if (retVal == NULL) {
        printf ("Error invoking function [%s]\n", ruleActionFnName);
        Py_DECREF(pyFunc);
        Py_DECREF(retVal);
        Py_Finalize();
        return -1;
   }
   Py_DECREF(pyFunc);

   PyArg_Parse(retVal, "i", &iRetVal);
   Py_DECREF(retVal);


   Py_Finalize();
   printf("Closed Python \n");
   return 0;
}