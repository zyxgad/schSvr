
package kpnmsqlcli

import (
	time    "time"
	unsafe  "unsafe"
	// reflect "reflect"
	// strings "strings"
)

type Map map[string]interface{}
type WhereValue struct{
	Key   string
	Cond  string
	Value interface{}
	Next  string
}
type WhereMap []WhereValue
/*
 * use it like this:
 * WhereMap{
 *   {"name", "=", "steve", "AND"},
 *   {"age", ">=", 18, ""},
 * }
 * will format to:
 *   "WHERE `name` = ?"
 */

type TypeMap map[string]string

type DateTime time.Time


func WhereValueToArr(wval WhereValue)(arr []interface{}){
	arr = make([]interface{}, 4)
	arr[0] = wval.Key
	arr[1] = wval.Cond
	arr[2] = wval.Value
	arr[3] = wval.Next
	return arr
}

func ArrToWhereValue(arr []interface{})(wval WhereValue){
	return WhereValue{
		Key: arr[0].(string),
		Cond: arr[1].(string),
		Value: arr[2],
		Next: arr[3].(string),
	}
}

func WhereMapToArr(wmap WhereMap)(arr [][]interface{}){
	arr = make([][]interface{}, len(wmap))
	for i, _ := range wmap {
		arr[i] = WhereValueToArr(wmap[i])
	}
	return arr
}

func ArrToWhereMap(arr [][]interface{})(wmap WhereMap){
	wmap = make(WhereMap, len(arr))
	for i, _ := range arr {
		wmap[i] = ArrToWhereValue(arr[i])
	}
	return wmap
}

func getMapKValues(obj Map)(keys []string, values []interface{}){
	if obj == nil {
		return nil, nil
	}
	keys = make([]string, 0, len(obj))
	values = make([]interface{}, 0, len(obj))
	for k, v := range obj {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func mapArrToMapArr(amap []map[string]interface{})(marr []Map){
	marr = *(*[]Map)(unsafe.Pointer(&amap))
	return marr
}

func FormatDatetime(t time.Time)(str string){
	return t.UTC().Format(FMT_Datetime)
}

func ParseDatetime(str string)(t time.Time){
	var err error
	t, err = time.Parse(FMT_Datetime, str)
	if err != nil {
		return time.Time{}
	}
	return t
}
