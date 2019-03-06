package main

import (
	"encoding/json"
	"github.com/project-flogo/rules/common/model"
)

//func TuplesFromJsonMap(tupleJsonStr string) map[model.TupleType]model.Tuple {
//	jsonMap := make(map[model.TupleType]interface{})
//	json.Unmarshal([]byte(tupleJsonStr), &jsonMap)
//	tupleMap := map[model.TupleType]model.Tuple{}
//
//	for k, v := range jsonMap {
//		tupleType := model.TupleType(k)
//		tupMap := v.(map[string]interface{})["Tuples"].(map[string]interface{})
//		tuple, _ := model.NewTuple(tupleType, tupMap)
//		tupleMap [tupleType] = tuple
//	}
//	return tupleMap
//}

func TupleFromJsonStr(tupleJsonStr string) model.Tuple {
	jsonMap := make(map[model.TupleType]interface{})
	json.Unmarshal([]byte(tupleJsonStr), &jsonMap)
	tupleType := jsonMap["TupleType"].(string)
	tupleProps := jsonMap["Tuples"].(map[string]interface{})
	tuple, _ := model.NewTuple(model.TupleType(tupleType), tupleProps)
	return tuple
}

//func TupleFromJson(tupleJsonStr string) map[model.TupleType]model.Tuple {
//	jsonMap := make(map[model.TupleType]interface{})
//	json.Unmarshal([]byte(tupleJsonStr), &jsonMap)
//	tupleMap := map[model.TupleType]model.Tuple{}
//
//	for k, v := range jsonMap {
//		tupleType := model.TupleType(k)
//		tupMap := v.(map[string]interface{})["Tuples"].(map[string]interface{})
//		tuple, _ := model.NewTuple(tupleType, tupMap)
//		tupleMap [tupleType] = tuple
//	}
//	return tupleMap
//}