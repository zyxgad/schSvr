
package kpnmsql

import (
	net    "net"

	util   "github.com/zyxgad/go-util/util"
	utcp   "github.com/zyxgad/go-util/tcp"
)

var (
	server *utcp.Server
)


var serverPacketCall = func(conn net.Conn, pkt utcp.Packet){
	var response = util.JsonType{
		"status": "",
	}
	defer func(){
		logger.Debugln("response:", response)
		packet := utcp.NewPacket(nil)
		packet.SetJson(response)
		_, err := packet.WriteTo(conn)
		if err != nil {
			logger.Fatalln("Send data error:", err)
		}
	}()
	defer util.RecoverErr(func(err interface{}){
		response["status"] = "error"
		switch e := err.(type) {
		case *util.AssertErr:
			response["errorMessage"] = e.Msg()
		case error:
			response["errorMessage"] = e.Error()
		default:
			response["errorMessage"] = util.JoinObjStr("Sql error:", e)
		}
		logger.Errorln(response["errorMessage"])
		logger.Errorln(util.GetStack(0))
	})

	data, ok := pkt.ParseJson()

	if !ok {
		return
	}

	logger.Debugln("data:", data)

	table := NewSqlTable(util.JsonToString(data["table"]))
	mode := util.JsonToByte(data["mode"])
	switch mode {
	case SQL_INSERT_CODE:
		obj := (Map)(util.JsonToMap(data["obj"]))
		if table.SqlInsert(obj) {
			response["status"] = "ok"
		}else{
			response["status"] = "failed"
		}
	case SQL_UPDATE_CODE:
		obj := (Map)(util.JsonToMap(data["obj"]))
		wmap := ArrToWhereMap(util.JsonToDoubleArr(data["where"]))
		if table.SqlUpdate(obj, wmap) {
			response["status"] = "ok"
		}else{
			response["status"] = "failed"
		}
	case SQL_DELETE_CODE:
		wmap := ArrToWhereMap(util.JsonToDoubleArr(data["where"]))
		if table.SqlDelete(wmap) {
			response["status"] = "ok"
		}else{
			response["status"] = "failed"
		}
	case SQL_SEARCH_CODE:
		obj := (TypeMap)(util.JsonToMapString(data["obj"]))
		wmap := ArrToWhereMap(util.JsonToDoubleArr(data["where"]))
		var limit, offset uint = 0, 0
		if limit0, ok := data["limit"]; ok { limit = uint(util.JsonToUint32(limit0)) }
		if offset0, ok := data["offset"]; ok { offset = uint(util.JsonToUint32(offset0)) }
		if lines := table.SqlSearch(obj, wmap, offset, limit); lines != nil {
			response["status"] = "ok"
			response["data"] = lines
		}else{
			response["status"] = "failed"
		}
	default:
		response["status"] = "failed"
		response["errorMessage"] = util.JoinObjStr("Unknow sql mode:", mode)
	}
}


type serverSource int

func (serverSource)Init(){
	server = utcp.NewServer(MAX_CONN, HOST + ":" + PORT)
}

func StartServer(){
	logger.Warnln("Sql server starting...")
	server.ServerStart(serverPacketCall)
}
