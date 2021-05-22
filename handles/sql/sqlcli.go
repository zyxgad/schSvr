
package kpnmsqlcli

import (
	sync "sync"

	utcp "github.com/zyxgad/go-util/tcp"
	util "github.com/zyxgad/go-util/util"
)

const CONN_NUM = 16

type Connect struct{
	using bool
	conn *utcp.Client
}

func (conni *Connect)send(data util.JsonType)(err error){
	if !conni.conn.IsConn() {
		conni.conn.Conn()
	}

	pkt := utcp.NewPacket(nil)
	pkt.SetJson(data)

	return conni.conn.SendPkt(pkt)
}

func (conni *Connect)recv()(data util.JsonType, err error){
	if !conni.conn.IsConn() {
		conni.conn.Conn()
	}

	var pkt utcp.Packet
	pkt, err = conni.conn.RecvPkt()
	if err != nil {
		return nil, err
	}
	data, _ = pkt.ParseJson()
	return data, nil
}


type Operator struct{
	name string
	lock sync.Locker

	conn_pool []*Connect
	free_chan chan bool
}

func NewOperator(name string)(*Operator){
	var conn_pool []*Connect = make([]*Connect, CONN_NUM)

	for i := 0; i < CONN_NUM ;i++ {
		conn_pool[i] = &Connect{
			using: false,
			conn: utcp.NewClient(SQL_HOST + ":" + SQL_PORT),
		}
	}

	return &Operator{
		name: name,
		lock: new(sync.Mutex),
		conn_pool: conn_pool,
		free_chan: make(chan bool),
	}
}

// func (opr *Operator)IsConn()(bool){
// 	return opr != nil && opr.cli.IsConn()
// }

// func (opr *Operator)Conn()(err error){
// 	return opr.cli.Conn()
// }

// func (opr *Operator)Close()(err error){
// 	return opr.cli.Close()
// }

// func (opr *Operator)send(data util.JsonType)(err error){
// 	if !opr.cli.IsConn() {
// 		opr.Conn()
// 	}

// 	pkt := utcp.NewPacket(nil)
// 	pkt.SetJson(data)

// 	return opr.cli.SendPkt(pkt)
// }

// func (opr *Operator)recv()(data util.JsonType, err error){
// 	if !opr.cli.IsConn() {
// 		opr.Conn()
// 	}

// 	var pkt utcp.Packet
// 	pkt, err = opr.cli.RecvPkt()
// 	if err != nil {
// 		return nil, err
// 	}
// 	data, _ = pkt.ParseJson()
// 	return data, nil
// }

func (opr *Operator)getconn()(*Connect){
	opr.lock.Lock()
	defer opr.lock.Unlock()
	for _, conni := range opr.conn_pool {
		if !conni.using {
			conni.using = true
			return conni
		}
	}
	return nil
}

func (opr *Operator)Getconn()(conni *Connect){
	for{
		conni = opr.getconn()
		if conni != nil {
			return conni
		}
		opr.free_chan <- true
	}
	return nil
}

func (opr *Operator)Freeconn(conni *Connect){
	conni.using = false
	select{ case <-opr.free_chan:default: }
}

func (opr *Operator)NewTable(name string)(tb SqlTable){
	return &Table{
		name: name,
		opr: opr,
	}
}

// func (opr *Operator)Lock(){
// 	opr.lock.Lock()
// }

// func (opr *Operator)Unlock(){
// 	opr.lock.Unlock()
// }


