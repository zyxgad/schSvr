
package kpnmsqlcli

const (
	_ byte = iota
	SQL_INSERT_CODE
	SQL_UPDATE_CODE
	SQL_DELETE_CODE
	SQL_SEARCH_CODE
)

const (
	TYPE_Bool    = "bool"
	TYPE_Int8    = "int8"
	TYPE_Int16   = "int16"
	TYPE_Int     = "int"
	TYPE_Int32   = "int32"
	TYPE_Int64   = "int64"
	TYPE_Uint8   = "uint8"
	TYPE_Uint16  = "uint16"
	TYPE_Uint    = "uint"
	TYPE_Uint32  = "uint32"
	TYPE_Uint64  = "uint64"
	TYPE_Float32 = "float32"
	TYPE_Float64 = "float64"
	TYPE_String  = "string"

	TYPE_Datetime = "DateTime"
)

const (
	FMT_Datetime = "2006-01-02 15:04:05"
)
