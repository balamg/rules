package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/rules/common"
)

func TestOne(t *testing.T) {

	td1 := TupleDescriptor{}
	td1.Name = "a"
	td1.TTLInSeconds = 10
	td1p1 := TuplePropertyDescriptor{}
	td1p1.Name = "p1"
	td1p1.KeyIndex = 3
	td1p1.PropType = data.TypeFloat64
	td1p2 := TuplePropertyDescriptor{}
	td1p2.Name = "p2"
	td1p2.KeyIndex = 31
	td1p2.PropType = data.TypeString

	td1.Props = []TuplePropertyDescriptor{td1p1, td1p2}

	str, _ := json.Marshal(&td1)
	t.Logf("succes %s\n", str)

	tpdx := TupleDescriptor{}
	tpdx.TTLInSeconds = -1
	json.Unmarshal([]byte(str), &tpdx)

	str1, _ := json.Marshal(&tpdx)
	t.Logf("succes %s\n", str1)

}

func TestTwo(t *testing.T) {
	tupleDescAbsFileNm := common.GetAbsPathForResource("src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	t.Logf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	RegisterTupleDescriptors(tupleDescriptor)

}

func Test_TupleJson (t *testing.T) {
	tupleDescAbsFileNm := common.GetAbsPathForResource("src/github.com/project-flogo/rules/pyruletest/pyrulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	//fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	var mp = map[string]interface{} {
		"name" : "Bala",
		"age" : 48,
		"gender" : "Male",
		"salary" : 100.1212,
	}
	t1, err := NewTuple("n1", mp)

	var mp2 = map[string]interface{} {
		"name" : "Supriya",
	}
	t2, err := NewTuple("n2", mp2)


	var tuples = map [TupleType]Tuple {
		TupleType("n1") : t1,
		TupleType("n2") : t2,
	}

	//tuples := []Tuple {t1, t2}

	b, err := json.Marshal(tuples)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return;
	}
	fmt.Println("**"+string(b)+"**")



	unmarshalled := make(map[TupleType]interface{})

	json.Unmarshal(b, &unmarshalled)
	ts := TuplesFromJsonMap(unmarshalled)
	fmt.Printf("done.. %v\n", ts)

	x := "{\"TupleType\": \"n1\", \"Tuples\": {\"salary\": 100.1, \"gender\": \"Male\", \"age\": 48, \"name\": \"Bala\"}}"
	fmt.Printf("&&%s&&\n", x)
	ts1 := TupleFromJsonStr(x)
	fmt.Printf("done.. %v\n", ts1)

	y := "[\"n1\", \"n2\"]"

	by := []byte (y)
	strarr := []string{}
	json.Unmarshal(by,&strarr)

	fmt.Printf("done")

}

func TuplesFromJsonMap(tupleJson map[TupleType]interface{}) map[TupleType]Tuple {
	tupleMap := map[TupleType]Tuple{}
	for k, v := range tupleJson {
		tupleType := TupleType(k)
		tupMap := v.(map[string]interface{})["Tuples"].(map[string]interface{})
		tuple, _ := NewTuple(tupleType, tupMap)
		tupleMap [tupleType] = tuple
	}
	return tupleMap
}

//func TuplesFromJsonMapStr(tupleJsonStr string) map[TupleType]Tuple {
//	jsonMap := make(map[TupleType]interface{})
//	json.Unmarshal([]byte(tupleJsonStr), &jsonMap)
//	tupleMap := map[TupleType]Tuple{}
//
//	for k, v := range jsonMap {
//		tupleType := TupleType(k)
//		tupMap := v.(map[string]interface{})["Tuples"].(map[string]interface{})
//		tuple, _ := NewTuple(tupleType, tupMap)
//		tupleMap [tupleType] = tuple
//	}
//	return tupleMap
//}

//func TupleFromJson(tupleJson map[TupleType]interface{}) Tuple {
//	tupleType := tupleJson["TupleType"].(string)
//	tupleProps := tupleJson["Tuples"].(map[string]interface{})
//	tuple, _ := NewTuple(TupleType(tupleType), tupleProps)
//
//	return tuple
//}

func TupleFromJsonStr(tupleJsonStr string) Tuple {
	jsonMap := make(map[TupleType]interface{})
	json.Unmarshal([]byte(tupleJsonStr), &jsonMap)
	tupleType := jsonMap["TupleType"].(string)
	tupleProps := jsonMap["Tuples"].(map[string]interface{})
	tuple, _ := NewTuple(TupleType(tupleType), tupleProps)
	return tuple
}
