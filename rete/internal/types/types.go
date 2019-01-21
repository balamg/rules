package types

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/common"
	"context"
)

type Network interface {
	common.Network
	GetPrefix() string
	GetIdGenService() IdGen
	GetJtService() JtService
	GetHandleService() HandleService
	GetJtRefService() JtRefsService
	GetTupleStore() model.TupleStore
}

type JoinTable interface {
	NwElemId
	GetName() string
	GetRule() model.Rule

	AddRow(ctx context.Context, handles []ReteHandle) JoinTableRow
	RemoveRow(ctx context.Context, rowID int) JoinTableRow
	GetRow(rowID int) JoinTableRow
	GetRowIterator() JointableRowIterator

	GetRowCount() int
	RemoveAllRows(ctx context.Context) //used when join table needs to be deleted
}

type JoinTableRow interface {
	NwElemId
	GetHandles() []ReteHandle
}

type ReteHandle interface {
	NwElemId
	SetTuple(tuple model.Tuple)
	GetTuple() model.Tuple
	GetTupleKey() model.TupleKey
}

type JtRefsService interface {
	NwService
	AddEntry(ctx context.Context, handle ReteHandle, jtName string, rowID int)
	RemoveRowEntry(ctx context.Context, handle ReteHandle, jtName string, rowID int)
	RemoveTableEntry(ctx context.Context, handle ReteHandle, jtName string)
	RemoveEntry(ctx context.Context, handle ReteHandle, jtName string)
	GetTableIterator(handle ReteHandle) JointableIterator
	GetRowIterator(handle ReteHandle, jtName string) JointableRowIterator
}

type JtService interface {
	NwService
	GetOrCreateJoinTable(ctx context.Context, nw Network, rule model.Rule, identifiers []model.TupleType, name string) JoinTable
	GetJoinTable(name string) JoinTable
}

type HandleService interface {
	NwService
	RemoveHandle(ctx context.Context, tuple model.Tuple) ReteHandle
	GetHandle(tuple model.Tuple) ReteHandle
	GetHandleByKey(key model.TupleKey) ReteHandle
	GetOrCreateHandle(ctx context.Context, nw Network, tuple model.Tuple) ReteHandle
}

type IdGen interface {
	NwService
	GetMaxID() int
	GetNextID() int
}

type JointableIterator interface {
	HasNext() bool
	Next() JoinTable
	Remove(ctx context.Context)
}

type JointableRowIterator interface {
	HasNext() bool
	Next() JoinTableRow
	Remove(ctx context.Context)
}
