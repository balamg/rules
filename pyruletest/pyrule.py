from ctypes import *

lib = cdll.LoadLibrary("./capi.so")

class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]
# describe and invoke Add()
lib.RegisterTupleDescriptors.argtypes = [GoString]
#lib.Add.restype = c_longlong
#with open('/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json', 'r') as myfile:
  #data = myfile.read()
data = open('/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json', 'r').read()
print(data)
gostr = GoString(data, len(data))
a = lib.RegisterTupleDescriptors(gostr)
