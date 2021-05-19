
package kpnmsqlcli

import(
	util "github.com/zyxgad/go-util/util"
	utcp "github.com/zyxgad/go-util/tcp"
)

type SqlTable interface{
	Name()(string)

	SqlInsert(obj Map)(err error)
	// Insert(obj SqlObjType)(err error)

	SqlDelete(cond WhereMap)(err error)

	SqlUpdate(obj Map, cond WhereMap)(err error)
	// Update(obj SqlObjType, cond WhereMap)(err error)

	SqlSearch(obj TypeMap, cond WhereMap, limit uint)(lines []Map, err error)
	SqlSearchOff(obj TypeMap, cond WhereMap, offset uint, limit uint)(lines []Map, err error)
	// Search(obj SqlObjType, cond WhereMap, limit uint)(lines []Map)

	HasData(cond WhereMap)(ok bool)
	DataCount(cond WhereMap)(leng uint64)
}

type Table struct{
	SqlTable

	name string
	opr *Operator
	conn *utcp.Client
}

func (tb *Table)Name()(string){
	return tb.name
}

func (tb *Table)sqlInsert(obj Map)(err error){
	var res util.JsonType
	res, err = tb.sendAndRecv(util.JsonType{
		"table": tb.name,
		"mode": SQL_INSERT_CODE,
		"obj": obj,
	})
	if err != nil {
		return err
	}
	if util.JsonToString(res["status"]) == "ok" {
		return nil
	}
	if util.JsonToString(res["status"]) == "error" {
		return util.NewErr(util.JoinObjStr("Sql insert error:", util.JsonToString(res["errorMessage"])), nil)
	}
	return util.NewErr("Sql insert server failed", nil)
}

func (tb *Table)SqlInsert(obj Map)(err error){
	util.Assert(obj != nil && len(obj) > 0, "Obj can't be nil and it's length must greater than 0")

	return tb.sqlInsert(obj)
}

// func (tb *Table)Insert(obj SqlObjType)(ok bool){
// 	util.Assert(obj != nil, "Obj can't be nil")

// 	keys := getStructTags(obj, true)
// 	values := getStructValues(obj)
// 	return tb.sqlInsert(keys, values)
// }

func (tb *Table)sqlDelete(cond WhereMap)(err error){
	var res util.JsonType
	res, err = tb.sendAndRecv(util.JsonType{
		"table": tb.name,
		"mode": SQL_DELETE_CODE,
		"where": WhereMapToArr(cond),
	})
	if err != nil {
		return err
	}
	if util.JsonToString(res["status"]) == "ok" {
		return nil
	}
	if util.JsonToString(res["status"]) == "error" {
		return util.NewErr(util.JoinObjStr("Sql delete error:", util.JsonToString(res["errorMessage"])), nil)
	}
	return util.NewErr("Sql delete server failed", nil)
}

func (tb *Table)SqlDelete(cond WhereMap)(err error){
	return tb.sqlDelete(cond)
}

func (tb *Table)sqlUpdate(obj Map, cond WhereMap)(err error){
	var res util.JsonType
	res, err = tb.sendAndRecv(util.JsonType{
		"table": tb.name,
		"mode": SQL_UPDATE_CODE,
		"obj": obj,
		"where": WhereMapToArr(cond),
	})
	if err != nil {
		return err
	}
	if util.JsonToString(res["status"]) == "ok" {
		return nil
	}
	if util.JsonToString(res["status"]) == "error" {
		return util.NewErr(util.JoinObjStr("Sql update error:", util.JsonToString(res["errorMessage"])), nil)
	}
	return util.NewErr("Sql update server failed", nil)
}

func (tb *Table)SqlUpdate(obj Map, cond WhereMap)(err error){
	util.Assert(obj != nil && len(obj) > 0, "Obj can't be nil and it's length must greater than 0")

	if len(cond) == 0{
		return nil
	}

	return tb.sqlUpdate(obj, cond)
}

// func (tb *Table)Update(obj SqlObjType, cond WhereMap)(err error){
// 	util.Assert(obj != nil, "Obj can't be nil")

// 	if len(cond) == 0{
// 		return false
// 	}

// 	keys := getStructTags(obj, true)
// 	values := getStructValues(obj)
// 	return tb.sqlUpdate(keys, values, cond)
// }

func (tb *Table)sqlSearch(obj TypeMap, cond WhereMap, offset uint, limit uint)(lines []Map, err error){
	var res util.JsonType
	res, err = tb.sendAndRecv(util.JsonType{
		"table": tb.name,
		"mode": SQL_SEARCH_CODE,
		"obj": obj,
		"where": WhereMapToArr(cond),
		"offset": offset,
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}
	if util.JsonToString(res["status"]) == "ok" {
		lines = mapArrToMapArr(util.JsonToArrMap(res["data"]))
		return lines, nil
	}
	if util.JsonToString(res["status"]) == "error" {
		return nil, util.NewErr(util.JoinObjStr("Sql search error:", util.JsonToString(res["errorMessage"])), nil)
	}
	return nil, util.NewErr("Sql search server failed", nil)
}

func (tb *Table)SqlSearch(obj TypeMap, cond WhereMap, limit uint)(lines []Map, err error){
	util.Assert(obj != nil, "Obj can't be nil")
	return tb.sqlSearch(obj, cond, 0, limit)
}

func (tb *Table)SqlSearchOff(obj TypeMap, cond WhereMap, offset uint, limit uint)(lines []Map, err error){
	util.Assert(obj != nil, "Obj can't be nil")
	return tb.sqlSearch(obj, cond, offset, limit)
}

// func (tb *Table)Search(obj SqlObjType, cond WhereMap, limit uint)(lines []Map){
// 	util.Assert(obj != nil, "Obj can't be nil")
// 	if len(cond) == 0{
// 		return nil
// 	}

// 	tmap := getObjTypeMap(obj)

// 	return tb.SqlSearch(tmap, cond, limit)
// }

func (tb *Table)HasData(cond WhereMap)(ok bool){
	ls, err := tb.SqlSearch(TypeMap{"\\1": TYPE_Int}, cond, 1)
	if err != nil || len(ls) == 0 {
		return false
	}
	return true
}

func (tb *Table)DataCount(cond WhereMap)(leng uint64){
	ls, err := tb.SqlSearch(TypeMap{"\\1": TYPE_Int}, cond, 0)
	if err != nil || ls == nil {
		return 0
	}
	return (uint64)(len(ls))
}

func (tb *Table)sendAndRecv(data util.JsonType)(res util.JsonType, err error){
	tb.opr.Lock()
	defer tb.opr.Unlock()

	err = tb.opr.send(data)
	if err != nil {
		return nil, err
	}

	res, err = tb.opr.recv()
	if err != nil {
		return nil, err
	}

	return res, err
}


