package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
)

const (
	RTC_ISNEW = 0x01
	RTC_ISMODIFIED = 0x01 << 2
	RTC_ISDELETED = 0x01 << 3
)

//Holds a tuple reference and related state
type reteHandle interface {
	setTuple(tuple model.Tuple)
	getTuple() model.Tuple
	addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable)
	removeJoinTableRowRefs(changedProps map[string]bool)
	removeJoinTable(joinTableVar joinTable)

	setNewlyAsserted()
	setModified()
	setRetracted()

	isNewlyAsserted() bool
	isModified() bool
	isRetracted() bool

}

type handleImpl struct {
	tuple         model.Tuple
	tablesAndRows map[joinTable]*list.List

	rtcBits		  uint8
}

func (hdl *handleImpl) setTuple(tuple model.Tuple) {
	hdl.tuple = tuple
}

func (hdl *handleImpl) initHandleImpl() {
	hdl.tablesAndRows = make(map[joinTable]*list.List)
	hdl.rtcBits = 0x00
}

func (hdl *handleImpl) getTuple() model.Tuple {
	return hdl.tuple
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) reteHandle {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(tuple)
}

func (hdl *handleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {

	rowsForJoinTable := hdl.tablesAndRows[joinTableVar]
	if rowsForJoinTable == nil {
		rowsForJoinTable = list.New()
		hdl.tablesAndRows[joinTableVar] = rowsForJoinTable
	}
	rowsForJoinTable.PushBack(joinTableRowVar)

}

func (hdl *handleImpl) removeJoinTableRowRefs(changedProps map[string]bool) {

	tuple := hdl.tuple
	alias := tuple.GetTupleType()

	emptyJoinTables := list.New()

	for joinTable, listOfRows := range hdl.tablesAndRows {

		toDelete := false
		if changedProps != nil {
			rule := joinTable.getRule()
			depProps, found := rule.GetDeps()[alias]
			if found { // rule depends on this type
				for changedProp := range changedProps {
					_, foundProp := depProps[changedProp]
					if foundProp {
						toDelete = true
						break
					}
				}
			}
		} else {
			toDelete = true
		}

		if !toDelete {
			continue
		}

		for e := listOfRows.Front(); e != nil; e = e.Next() {
			row := e.Value.(joinTableRow)
			joinTable.removeRow(row)
		}
		if joinTable.len() == 0 {
			emptyJoinTables.PushBack(joinTable)
		}
	}

	for e := emptyJoinTables.Front(); e != nil; e = e.Next() {
		emptyJoinTable := e.Value.(joinTable)
		delete(hdl.tablesAndRows, emptyJoinTable)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *handleImpl) removeJoinTable(joinTableVar joinTable) {
	_, ok := hdl.tablesAndRows[joinTableVar]
	if ok {
		delete(hdl.tablesAndRows, joinTableVar)
	}
}
func (hdl *handleImpl) setNewlyAsserted() {

}
func (hdl *handleImpl) setModified() {

}
func (hdl *handleImpl) setRetracted() {

}

func (hdl *handleImpl) isNewlyAsserted() bool {

}
func (hdl *handleImpl) isModified() bool {

}
func (hdl *handleImpl) isRetracted() bool {

}