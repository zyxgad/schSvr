
package kpnmsql

import (
	// kutil "github.com/zyxgad/KpnmServer/util"
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
 * MakeWhere(WhereMap{
 *   {"name", "=", "steve", "AND"},
 *   {"age", ">=", 18, ""},
 * })
 * will return:
 *   "WHERE `name` = ?"
 */

type TypeMap map[string]string


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


func FormatWhere(argv WhereMap)(str string, values[]interface{}){
	values = make([]interface{}, 0, len(argv))
	if len(argv) == 0 {
		return "", make([]interface{}, 0, 0)
	}


	str = "WHERE "
	for _, v := range argv {
		values = append(values, v.Value)
		str += "`" + v.Key + "` " + v.Cond + " ? " + v.Next
	}
	str = str[:len(str) - 1 - len(argv[len(argv) - 1].Next)]

	return str, values
}

func FormatValue(tags []string)(str string){
	if len(tags) == 0 {
		return ""
	}

	seats := ""
	str = ""
	for _, v := range tags {
		str += "`" + v + "`,"
		seats += "?,"
	}
	str = "(" + str[:len(str)-1] + ") VALUES (" + seats[:len(seats)-1] + ")"
	return str
}

func FormatSet(tags []string)(str string){
	if len(tags) == 0 {
		return ""
	}

	str = "SET "
	for _, v := range tags {
		str += "`" + v + "` = ?,"
	}
	str = str[:len(str)-1]
	return str
}

func FormatKeyList(tags []string)(str string){
	if len(tags) == 0 {
		return ""
	}

	str = ""
	for _, v := range tags {
		if v[0] == '\\' {
			str += v[1:] + ","
		}else{
			str += "`" + v + "`,"
		}
	}
	str = str[:len(str)-1]
	return str
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

func SplitKeyType(tpm TypeMap)(keys []string, types []string){
	if len(tpm) == 0 {
		return nil, nil
	}

	keys = make([]string, 0, len(tpm))
	types = make([]string, 0, len(tpm))

	for k, t := range tpm {
		keys = append(keys, k)
		types = append(types, t)
	}

	return keys, types
}

func BindKeyValue(keys []string, values []interface{})(obj Map){
	if len(keys) == 0 || len(keys) != len(values) {
		return nil
	}

	obj = make(Map)

	for i, _ := range keys {
		obj[keys[i]] = values[i]
	}

	return obj
}

func makeScanRowByTypes(types []string)(row []interface{}){
	row = make([]interface{}, 0, len(types))
	for _, t := range types {
		switch t {
		case TYPE_Bool:
			row = append(row, new(bool))
		case TYPE_Int8:
			row = append(row, new(int8))
		case TYPE_Int16:
			row = append(row, new(int16))
		case TYPE_Int:
			fallthrough
		case TYPE_Int32:
			row = append(row, new(int32))
		case TYPE_Int64:
			row = append(row, new(int64))
		case TYPE_Uint8:
			row = append(row, new(uint8))
		case TYPE_Uint16:
			row = append(row, new(uint16))
		case TYPE_Uint:
			fallthrough
		case TYPE_Uint32:
			row = append(row, new(uint32))
		case TYPE_Uint64:
			row = append(row, new(uint64))
		case TYPE_Float32:
			row = append(row, new(float32))
		case TYPE_Float64:
			row = append(row, new(float64))
		case TYPE_String:
			row = append(row, new(string))
		case "time.Time":
			fallthrough
		case TYPE_Datetime:
			row = append(row, new(string))
		default:
			panic("Unknow type " + t)
		}
	}
	return row
}
