
package kpnmlogcli

import (
	time "time"

	util  "github.com/zyxgad/go-util/util"
	utcp  "github.com/zyxgad/go-util/tcp"
	kutil "github.com/zyxgad/schSvr/util"
)


type Logger struct{
	name string
	cli *utcp.Client
}

func NewLogger(name string)(*Logger){
	return &Logger{
		name: name,
		cli: utcp.NewClient(LOGGER_HOST + ":" + LOGGER_PORT),
	}
}

func (lgr *Logger)IsConn()(bool){
	return lgr != nil && lgr.cli.IsConn()
}

func (lgr *Logger)Conn()(err error){
	return lgr.cli.Conn()
}

func (lgr *Logger)Close()(err error){
	return lgr.cli.Close()
}

func (lgr *Logger)send(msg string)(err error){
	pkt := utcp.NewPacket(nil)
	pkt.SetString(msg)

	return lgr.cli.SendPkt(pkt)
}

/************************************************************/
func (lgr *Logger)Debugln(obj ...interface{})(err error){
	if !kutil.IsDebug(){
		return nil
	}
	str := util.JoinObjStr(obj...)
	msg := util.FormatStr("[%s:DEBUG][%s]: %s", lgr.name, getFmtTime(true), str)
	return lgr.send(msg)
}

func (lgr *Logger)Debugf(format string, obj ...interface{})(err error){
	if !kutil.IsDebug(){
		return nil
	}
	str := util.FormatStr(format, obj...)
	msg := util.FormatStr("[%s:DEBUG][%s]: %s", lgr.name, getFmtTime(true), str)
	return lgr.send(msg)
}

func (lgr *Logger)Infoln(obj ...interface{})(err error){
	str := util.JoinObjStr(obj...)
	msg := util.FormatStr("[%s:INFO][%s]: %s", lgr.name, getFmtTime(false), str)
	return lgr.send(msg)
}

func (lgr *Logger)Infof(format string, obj ...interface{})(err error){
	str := util.FormatStr(format, obj...)
	msg := util.FormatStr("[%s:INFO][%s]: %s", lgr.name, getFmtTime(false), str)
	return lgr.send(msg)
}

func (lgr *Logger)Warnln(obj ...interface{})(err error){
	str := util.JoinObjStr(obj...)
	msg := util.FormatStr("[%s:WARN][%s]: %s", lgr.name, getFmtTime(false), str)
	return lgr.send(msg)
}

func (lgr *Logger)Warnf(format string, obj ...interface{})(err error){
	str := util.FormatStr(format, obj...)
	msg := util.FormatStr("[%s:WARN][%s]: %s", lgr.name, getFmtTime(false), str)
	return lgr.send(msg)
}

func (lgr *Logger)Errorln(obj ...interface{})(err error){
	str := util.JoinObjStr(obj...)
	msg := util.FormatStr("[%s:ERROR][%s]: %s", lgr.name, getFmtTime(true), str)
	return lgr.send(msg)
}

func (lgr *Logger)Errorf(format string, obj ...interface{})(err error){
	str := util.FormatStr(format, obj...)
	msg := util.FormatStr("[%s:ERROR][%s]: %s", lgr.name, getFmtTime(true), str)
	return lgr.send(msg)
}

func (lgr *Logger)Fatalln(obj ...interface{})(err error){
	str := util.JoinObjStr(obj...)
	msg := util.FormatStr("[%s:FATAL][%s]: %s", lgr.name, getFmtTime(true), str)
	return lgr.send(msg)
}

func (lgr *Logger)Fatalf(format string, obj ...interface{})(err error){
	str := util.FormatStr(format, obj...)
	msg := util.FormatStr("[%s:FATAL][%s]: %s", lgr.name, getFmtTime(true), str)
	return lgr.send(msg)
}
/************************************************************/


func getFmtTime(nano bool)(string){
	tmnow := time.Now()
	unixtm := tmnow.Unix()
	if !nano {
		return util.FormatStr("%d", unixtm)
	}
	nanotm := tmnow.UnixNano() % 1e9
	return util.FormatStr("%d.%d", unixtm, nanotm)
}
