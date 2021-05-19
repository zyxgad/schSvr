
package kpnmrobot

import (
	bytes   "bytes"

	captcha "github.com/dchest/captcha"
	ksql    "github.com/zyxgad/schSvr/handles/sql"
	util    "github.com/zyxgad/go-util/util"
)

const (
	CAPT_LENGTH = 6
	CAPT_WIDTH = 240
	CAPT_HEIGHT = 100
)

var (
	sqlCaptTable ksql.SqlTable
)

func NewCaptcha()(id string, imgdata string, err error){
	id = captcha.NewLen(CAPT_LENGTH)

	var imgbuf *bytes.Buffer = bytes.NewBuffer([]byte{})
	err = captcha.WriteImage(imgbuf, id, CAPT_WIDTH, CAPT_HEIGHT)
	if err != nil {
		return "", "", err
	}

	imgdata = "data:image/png;base64," + util.EncodeBase64(imgbuf.Bytes())
	return id, imgdata, nil
}

func VerifyCaptcha(id string, value string)(ok bool){
	return captcha.VerifyString(id, value)
}

func RemoveCaptcha(id string)(ok bool){
	return sqlCaptTable.SqlDelete(ksql.WhereMap{{"id", "=", id, ""}}) != nil
}

type sqlCaptStore struct{
	captcha.Store

	sqltb ksql.SqlTable
}

func (cst *sqlCaptStore)Set(id string, digits []byte){
	var digstr []byte = make([]byte, len(digits))
	for i, d := range digits {
		digstr[i] = '0' + d
	}
	cst.sqltb.SqlInsert(ksql.Map{
		"id": id,
		"value": (string)(digstr),
		"overtime": ksql.FormatDatetime(util.GetTimeAfter(util.TimeMin * 10)),
	})
}

func (cst *sqlCaptStore)Get(id string, clear bool)(digits []byte){
	lines, err := cst.sqltb.SqlSearch(ksql.TypeMap{ "value": ksql.TYPE_String }, ksql.WhereMap{{"id", "=", id, ""}}, 1)
	if err != nil || len(lines) != 1 {
		return nil
	}
	if clear {
		cst.sqltb.SqlDelete(ksql.WhereMap{{"id", "=", id, ""}})
	}
	value := util.JsonToString(lines[0]["value"])
	digits = make([]byte, len(value))
	for i, d := range ([]byte)(value) {
		if d < '0' || '9' < d {
			return nil
		}
		digits[i] = d - '0'
	}
	return digits
}

func InitCaptchaSqlTable(tb ksql.SqlTable){
	if tb == nil || sqlCaptTable != nil {
		panic("nil table or it is inited")
		return
	}
	sqlCaptTable = tb
	captcha.SetCustomStore(&sqlCaptStore{ sqltb: sqlCaptTable })
}
