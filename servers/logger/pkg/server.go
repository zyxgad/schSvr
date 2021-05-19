
package kpnmlogger

import (
	net    "net"

	utcp "github.com/zyxgad/go-util/tcp"
)

var (
	server *utcp.Server
)

var serverPacketCall = func(conn net.Conn, pkt utcp.Packet){
	msg, ok := pkt.ParseString()

	if !ok {
		return
	}
	_logger.Println(msg)
}


type serverSource int

func (serverSource)Init(){
	server = utcp.NewServer(MAX_CONN, HOST + ":" + PORT)
}

func StartServer(){
	server.ServerStart(serverPacketCall)
}
