package redis

import (
	"context"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
	"strconv"
)

type jtRefsServiceImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	types.NwServiceImpl

	//tablesAndRows map[string]map[string]map[int]int
}

func NewJoinTableRefsInHdlImpl(nw types.Network, config map[string]interface{}) types.JtRefsService {
	hdlJt := jtRefsServiceImpl{}
	hdlJt.Nw = nw
	//hdlJt.tablesAndRows = make(map[string]map[string]map[int]int)
	return &hdlJt
}

func (h *jtRefsServiceImpl) Init() {

}

func (h *jtRefsServiceImpl) AddEntry(ctx context.Context, handle types.ReteHandle, jtName string, rowID int) {

	//format: prefix:rtbls:tkey ==> {jtname=jtname, ...}
	//format: prefix:rrows:tkey:jtname ==> {rowid=rowid, ...}

	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl := redisutils.GetRedisHdl()
	valMap := make(map[string]interface{})
	valMap[jtName] = jtName
	hdl.HSetAll(key, valMap)

	rkey := h.Nw.GetPrefix() + ":rrows:" + handle.GetTupleKey().String() + ":" + jtName
	rowMap := make(map[string]interface{})
	rowIdStr := strconv.Itoa(rowID)
	rowMap[rowIdStr] = rowIdStr
	hdl.HSetAll(rkey, rowMap)
}

func (h *jtRefsServiceImpl) RemoveEntry(ctx context.Context, handle types.ReteHandle, jtName string) {
	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl := redisutils.GetRedisHdl()
	hdl.HDel(key, jtName)

	rkey := h.Nw.GetPrefix() + ":rrows:" + handle.GetTupleKey().String() + ":" + jtName
	hdl.Del(rkey)

}

func (h *jtRefsServiceImpl) RemoveRowEntry(ctx context.Context, handle types.ReteHandle, jtName string, rowID int) {
	rowKey := h.Nw.GetPrefix() + ":rrows:" + handle.GetTupleKey().String() + ":" + jtName
	hdl := redisutils.GetRedisHdl()
	rowIdStr := strconv.Itoa(rowID)
	hdl.HDel(rowKey, rowIdStr)

	//hkey := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	//hdl.HDel(hkey, jtName)
}

func (h *jtRefsServiceImpl) RemoveTableEntry(ctx context.Context, handle types.ReteHandle, jtName string) {

	//Delete all row entry refs for this table
	rowKey := h.Nw.GetPrefix() + ":rrows:" + handle.GetTupleKey().String() + ":" + jtName
	hdl := redisutils.GetRedisHdl()
	hdl.Del(rowKey)

	//Delete the table entry for this key
	hkey := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl.HDel(hkey, jtName)
}

func (h *jtRefsServiceImpl) GetTableIterator(handle types.ReteHandle) types.JointableIterator {
	ri := hdlRefsTableIteratorImpl{}
	ri.nw = h.Nw
	//format: prefix:rtbls:tkey ==> {jtname=jtname, ...}
	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl := redisutils.GetRedisHdl()
	ri.iter = hdl.GetMapIterator(key)
	return &ri
}

type hdlRefsTableIteratorImpl struct {
	iter *redisutils.MapIterator
	nw   types.Network
}

func (ri *hdlRefsTableIteratorImpl) HasNext() bool {
	return ri.iter.HasNext()
}

func (ri *hdlRefsTableIteratorImpl) Next() types.JoinTable {
	jtName, _ := ri.iter.Next()
	jT := ri.nw.GetJtService().GetJoinTable(jtName)
	return jT
}
func (ri *hdlRefsTableIteratorImpl) Remove(ctx context.Context) {
	ri.iter.Remove()
}

type hdlRefsRowIteratorImpl struct {
	key    string
	iter   *redisutils.MapIterator
	nw     types.Network
	jtName string
}

func (r *hdlRefsRowIteratorImpl) HasNext() bool {
	return r.iter.HasNext()
}

func (r *hdlRefsRowIteratorImpl) Next() types.JoinTableRow {
	rowIdStr, _ := r.iter.Next()
	rowID, _ := strconv.Atoi(rowIdStr)
	jT := r.nw.GetJtService().GetJoinTable(r.jtName)
	row := jT.GetRow(rowID)
	return row
}

func (r *hdlRefsRowIteratorImpl) Remove(ctx context.Context) {
	r.iter.Remove()
}

//format: prefix:rtbls:tkey ==> {jtname=jtname, ...}
//format: prefix:rrows:tkey:jtname ==> {rowid=rowid, ...}

func (h *jtRefsServiceImpl) GetRowIterator(handle types.ReteHandle, jtName string) types.JointableRowIterator {
	r := hdlRefsRowIteratorImpl{}
	r.nw = h.Nw
	r.jtName = jtName
	//ex: a:rrows:n1:a:b1:L_tbl
	r.key = h.Nw.GetPrefix() + ":rrows:" + handle.GetTupleKey().String() + ":" + jtName
	hdl := redisutils.GetRedisHdl()
	r.iter = hdl.GetMapIterator(r.key)
	return &r
}
