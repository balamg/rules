from cpy import *
import os

lib = cdll.LoadLibrary(os.path.dirname(__file__)+os.path.sep+"capi.so")

cCond = CFUNCTYPE(c_int, c_char_p, c_char_p, c_char_p)
ccB = cCond(PyConditionCb)
lib.registerConditionCb(ccB)

cActionCb = CFUNCTYPE(c_int, c_char_p, c_char_p)
cA = cActionCb(PyActionCb)
lib.registerActionCb(cA)