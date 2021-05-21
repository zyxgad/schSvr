
package kpnmsqlcli

import (
	sync "sync"

	utcp "github.com/zyxgad/go-util/tcp"
	util "github.com/zyxgad/go-util/util"
)


type Operator struct{
	name string
	cli *utcp.Client
	lock sync.Locker
}

func NewOperator(name string)(*Operator){
	return &Operator{
		name: name,
		cli: utcp.NewClient(SQL_HOST + ":" + SQL_PORT),
		lock: new(sync.Mutex),
	}
}

func (opr *Operator)IsConn()(bool){
	return opr != nil && opr.cli.IsConn()
}

func (opr *Operator)Conn()(err error){
	return opr.cli.Conn()
}

func (opr *Operator)Close()(err error){
	return opr.cli.Close()
}

func (opr *Operator)send(data util.JsonType)(err error){
	if !opr.cli.IsConn() {
		opr.Conn()
	}

	pkt := utcp.NewPacket(nil)
	pkt.SetJson(data)

	return opr.cli.SendPkt(pkt)
}

func (opr *Operator)recv()(data util.JsonType, err error){
	if !opr.cli.IsConn() {
		opr.Conn()
	}

	var pkt utcp.Packet
	pkt, err = opr.cli.RecvPkt()
	if err != nil {
		return nil, err
	}
	data, _ = pkt.ParseJson()
	return data, nil
}

func (opr *Operator)NewTable(name string)(tb SqlTable){
	var conn *utcp.Client = nil
	return &Table{
		name: name,
		opr: opr,
		conn: conn,
	}
}

func (opr *Operator)Lock(){
	opr.lock.Lock()
}

func (opr *Operator)Unlock(){
	opr.lock.Unlock()
}


