import json

class TupleDescriptor:
    def __init__(self, name, props, ttl=-1):
        self.name = name
        self.properties = props
        self.ttl = ttl

class TupleProperty:
    def __init__(self, a, b, c):
        self.name = a
        self.type = b
        self.pkindex = c

class TupleDescriptorEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, TupleProperty):
           return {"name" : o.name, "type" : o.type, "pk-index" : o.pkindex}
        else:
            return o.__dict__

class TupleDescriptorDecoder(json.JSONDecoder):
    def __init__(self, *args, **kwargs):
        json.JSONDecoder.__init__(self, object_hook=self.object_hook, *args, **kwargs)
        self.tps = list()
        self.tds = list()
    def object_hook(self, dct):
        if 'type' in dct:
            pkindex = "-1"
            if 'pk-index' in dct:
                pkindex = dct['pk-index']
            tp = TupleProperty (dct["name"], dct["type"], pkindex)
            self.tps.append(tp)
            return tp
        elif 'properties' in dct:
            ttl = -1
            if "ttl" in dct:
                ttl = dict["ttl"]
            td = TupleDescriptor (dct["name"], self.tps, ttl)
            self.tps = list()
            self.tds.append(td)
            return td

class Tuple:
    def __init__(self, tupleType, tuples):
        self.TupleType = tupleType
        self.Tuples = tuples

class TupleEncoder(json.JSONEncoder):
    def default(self, o):
        # if isinstance(o, Tuple):
        #     return {"name" : o.name, "type" : o.type, "pk-index" : o.pkindex}
        # else:
            return o.__dict__

## Not used, does not work properly
class TupleDecoder(json.JSONDecoder):
    def __init__(self, *args, **kwargs):
        json.JSONDecoder.__init__(self, object_hook=self.object_hook, *args, **kwargs)

    def object_hook(self, dct):
        if 'Tuples' in dct:
            tp = Tuple(dct["TupleType"], dct["Tuples"])
        else:
            return dct

def TupleFromJsonStr (tupleJsonStr):
    parsedJson = json.loads(tupleJsonStr)
    tupleMap = {}
    for outerKey, outerVal in parsedJson.iteritems():
        for innerKey, innerVal in outerVal.iteritems():
            if innerKey == "Tuples":
                tuple = Tuple (outerKey, innerVal)
                tupleMap[outerKey] = tuple
    return tupleMap

def TuplesFromJsonStr (tupleJsonStr):
    parsedJson = json.loads(tupleJsonStr)
    tupleMap = {}
    for outerKey, outerVal in parsedJson.iteritems():
        for innerKey, innerVal in outerVal.iteritems():
            if innerKey == "Tuples":
                tuple = Tuple (outerKey, innerVal)
                tupleMap[outerKey] = tuple
    return tupleMap

def TuplesToJsonStr (tuple):
    tupleJsonStr = json.dumps(tuple, cls=TupleEncoder)
    return tupleJsonStr

if __name__ == "__main__":
    #Construct a tuple descriptor from json and back
    tdJsonStr = open("/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json", 'r').read()
    print (tdJsonStr)
    tdFromJsonStr = json.loads(tdJsonStr, cls=TupleDescriptorDecoder)
    tdToJsonStr = json.dumps(tdFromJsonStr, cls=TupleDescriptorEncoder)
    print (tdToJsonStr)

    #Construct a tuple map from json and back
    tupleJsonStr = '{"n1":{"TupleType":"n1","Tuples":{"age":48,"gender":"Male","name":"Bala","salary":100.1212}},"n2":{"TupleType":"n2","Tuples":{"name":"Supriya"}}}'
    tuplesFromJson = TuplesFromJsonStr(tupleJsonStr)
    print (tuplesFromJson)
    tupleToJsonStr = json.dumps(tuplesFromJson, cls=TupleEncoder)
    print (tupleToJsonStr)

    #Construct a tuple from a map
    props = {}
    props['name'] = "Bala"
    props['age'] = 48
    props['gender'] = "Male"
    props['salary'] = 100.1
    tuple = Tuple("n1", props)
    #and serialize it
    tupleJsonStr = json.dumps(tuple, cls=TupleEncoder)

    print (tupleJsonStr)
