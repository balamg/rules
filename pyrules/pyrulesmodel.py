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



if __name__ == "__main__":
    jsn = open("/home/bala/go/src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json", 'r').read()
    print (jsn)
    fromJsn = json.loads(jsn, cls=TupleDescriptorDecoder)
    toJsn = json.dumps(fromJsn, cls=TupleDescriptorEncoder)
    print (toJsn)

    # tds1 = json.loads(y, cls=TupleDescriptorDecoder)
    # print ("done..\n", y)