package model

import "encoding/json"

//IdentifiersToString Take a slice of Identifiers and return a string representation
func IdentifiersToString(identifiers []TupleType) string {
	str := ""
	for _, idr := range identifiers {
		str += string(idr) + ", "
	}
	return str
}

// Contains returns true if an identifier exists in the identifier array
func Contains(identifiers []TupleType, toCheck TupleType) (bool, int) {
	for idx, id := range identifiers {
		if id == toCheck {
			return true, idx
		}
	}
	return false, -1
}

func TupleFromJsonStr(tupleJsonStr string) Tuple {
	jsonMap := make(map[TupleType]interface{})
	json.Unmarshal([]byte(tupleJsonStr), &jsonMap)
	tupleType := jsonMap["TupleType"].(string)
	tupleProps := jsonMap["Tuples"].(map[string]interface{})
	tuple, _ := NewTuple(TupleType(tupleType), tupleProps)
	return tuple
}
