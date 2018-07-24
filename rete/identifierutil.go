package rete

import "github.com/TIBCOSoftware/bego/common/model"

//These set operations are used in building the rete network. See Network.buildNetwork

//AppendIdentifiers ... Append identifiers from set2 to set1
func AppendIdentifiers(set1 []model.TupleTypeAlias, set2 []model.TupleTypeAlias) []model.TupleTypeAlias {
	union := []model.TupleTypeAlias{}
	union = append(union, set1...)
	union = append(union, set2...)
	return union
}

//ContainedByFirst ... true if second is a subset of first
func ContainedByFirst(first []model.TupleTypeAlias, second []model.TupleTypeAlias) bool {

	if len(second) == 0 {
		return true
	} else if len(first) == 0 {
		return false
	}
	for _, idFromSecond := range second {
		contains := false
		for _, idFromFirst := range first {
			if idFromSecond == idFromFirst {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true

}

//OtherTwoAreContainedByFirst ... true if second and third are part of first
func OtherTwoAreContainedByFirst(first []model.TupleTypeAlias, second []model.TupleTypeAlias, third []model.TupleTypeAlias) bool {
	return ContainedByFirst(first, second) && ContainedByFirst(first, third)
}

//UnionIdentifiers ... union of the first and second sets
func UnionIdentifiers(first []model.TupleTypeAlias, second []model.TupleTypeAlias) []model.TupleTypeAlias {
	union := []model.TupleTypeAlias{}
	union = append(union, first...)
	union = append(union, SecondMinusFirst(first, second)...)
	return union
}

//SecondMinusFirst ... returns elements in the second that arent in the first
func SecondMinusFirst(first []model.TupleTypeAlias, second []model.TupleTypeAlias) []model.TupleTypeAlias {
	minus := []model.TupleTypeAlias{}
outer:
	for _, idrSecond := range second {
		for _, idrFirst := range first {
			if idrSecond == idrFirst {
				continue outer
			}
		}
		minus = append(minus, idrSecond)
	}
	return minus
}

//IntersectionIdentifiers .. intersection of the two sets
func IntersectionIdentifiers(first []model.TupleTypeAlias, second []model.TupleTypeAlias) []model.TupleTypeAlias {
	intersect := []model.TupleTypeAlias{}
	for _, idrSecond := range second {
		for _, idrFirst := range first {
			if idrSecond == idrFirst {
				intersect = append(intersect, idrSecond)
			}
		}
	}
	return intersect
}

//EqualSets ... compare two identifiers based on their contents
func EqualSets(first []model.TupleTypeAlias, second []model.TupleTypeAlias) bool {
	return len(SecondMinusFirst(first, second)) == 0 && len(SecondMinusFirst(first, second)) == 0
}

//GetIndex ... return the index of thisIdr in identifiers
func GetIndex(identifiers []model.TupleTypeAlias, thisIdr model.TupleTypeAlias) int {
	for i, idr := range identifiers {
		if idr == thisIdr {
			return i
		}
		i++
	}
	return -1
}
