
package kpnmsqlcli

import (
	util "github.com/zyxgad/go-util/util"
)

type AutoCleanTable struct{
	SqlTable

	Table SqlTable
	Interval int64
	_last_clean_time int64
}

func NewAutoCleanTable(table SqlTable, interval int64)(ctb *AutoCleanTable){
	ctb = new(AutoCleanTable)
	ctb.Table = table
	if interval < 0 {
		interval = 0
	}
	ctb.Interval = interval
	ctb._last_clean_time = 0
	return ctb
}

func (ctb *AutoCleanTable)Clean(must bool){
	tmnow := util.GetTimeNow()
	tm := tmnow.Unix()
	if !must && ctb._last_clean_time + ctb.Interval > tm {
		return
	}
	ctb._last_clean_time = tm

	ctb.Table.SqlDelete(WhereMap{{"overtime", "<", FormatDatetime(tmnow), ""}})
}

func (ctb *AutoCleanTable)Name()(string){
	return ctb.Table.Name()
}

func (ctb *AutoCleanTable)SqlInsert(obj Map)(err error){
	ctb.Clean(false)
	return ctb.Table.SqlInsert(obj)
}

func (ctb *AutoCleanTable)SqlDelete(cond WhereMap)(err error){
	ctb.Clean(false)
	return ctb.Table.SqlDelete(cond)
}

func (ctb *AutoCleanTable)SqlUpdate(obj Map, cond WhereMap)(err error){
	ctb.Clean(false)
	cond = append(WhereMap{{"overtime", ">=", FormatDatetime(util.GetTimeNow()), "AND"}}, cond...)
	return ctb.Table.SqlUpdate(obj, cond)
}

func (ctb *AutoCleanTable)SqlSearch(obj TypeMap, cond WhereMap, limit uint)(lines []Map, err error){
	ctb.Clean(false)
	cond = append(WhereMap{{"overtime", ">=", FormatDatetime(util.GetTimeNow()), "AND"}}, cond...)
	return ctb.Table.SqlSearch(obj, cond, limit)
}

func (ctb *AutoCleanTable)SqlSearchOff(obj TypeMap, cond WhereMap, offset uint, limit uint)(lines []Map, err error){
	ctb.Clean(false)
	cond = append(WhereMap{{"overtime", ">=", FormatDatetime(util.GetTimeNow()), "AND"}}, cond...)
	return ctb.Table.SqlSearchOff(obj, cond, offset, limit)
}

func (ctb *AutoCleanTable)HasData(cond WhereMap)(ok bool){
	ctb.Clean(false)
	cond = append(WhereMap{{"overtime", ">=", FormatDatetime(util.GetTimeNow()), "AND"}}, cond...)
	return ctb.Table.HasData(cond)
}

func (ctb *AutoCleanTable)DataCount(cond WhereMap)(leng uint64){
	ctb.Clean(false)
	cond = append(WhereMap{{"overtime", ">=", FormatDatetime(util.GetTimeNow()), "AND"}}, cond...)
	return ctb.Table.DataCount(cond)
}
