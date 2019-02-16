from ctypes import *

lib = cdll.LoadLibrary("./capi.so")

lib.printHello()


def PyFunc (pystr):
    print ("OK it now prints pycb here %s\n", pystr)
    return 0;

PyFunc ("Bhalchandra")

#register a callback

CMPFUNC = CFUNCTYPE(c_int, c_char_p)

print ("heree..1")

cmp = CMPFUNC(PyFunc)

print ("heree..2")


lib.registerPyCb (cmp)

print ("heree..3")


lib.callPyCb ("Sup")

print ("heree..4")

class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

lib.RegisterTupleDescriptors.argtypes = [GoString]
#lib.Add.restype = c_longlong
#with open('/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json', 'r') as myfile:
#data = myfile.read()
data = open('/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json', 'r').read()
print(data)
gostr = GoString(data, len(data))
a = lib.RegisterTupleDescriptors(gostr)




