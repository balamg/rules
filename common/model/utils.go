package model

import (
	"context"

	"github.com/project-flogo/core/data/coerce"
)

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

func TuplesToMap(tuples map[TupleType]Tuple) map[string]interface{} {
	tupleMap := map[string]interface{}{}
	for ttype, tuple := range tuples {
		tm := map[string]interface{}{}
		tupleMap[string(ttype)] = tm
		for propNm, propVal := range tuple.GetMap() {
			if pMap, ok := propVal.(map[string]interface{}); ok { //map prop
				v, _ := coerce.ToObject(pMap) //deep copies
				tm[propNm] = v
			} else if pArr, ok := propVal.(map[string]interface{}); ok { //array prop
				v, _ := coerce.ToArray(pArr) //deep copies
				tm[propNm] = v
			} else {
				tm[propNm] = propVal
			}
		}
	}
	return tupleMap
}

func ToTuples(ctx context.Context, rs RuleSession, valueMap map[string]interface{}) ([]Tuple, error) {
	tuples := []Tuple{}
	//tupleMap := map[TupleType]Tuple{}
	for tType, tVs := range valueMap {
		vals := tVs.(map[string]interface{})
		tk, err := NewTupleKey(TupleType(tType), vals)
		if err != nil {
			return tuples, err
		}
		//search by its key, for tuple in the rule session. if found,
		//it is an update operation, else its a new tuple
		tuple := rs.GetAssertedTuple(tk)
		if tuple != nil {
			mTuple := tuple.(MutableTuple)
			//update the tuple to new values from the map
			err := mTuple.SetValues(ctx, vals)
			if err != nil {
				return tuples, err
			}
		} else {
			tuple, err = NewTuple(TupleType(tType), vals)
			if err != nil {
				return tuples, err
			}
		}
		tuples = append(tuples, tuple)
	}
	return tuples, nil
}
