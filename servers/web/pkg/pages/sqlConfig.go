
package kpnmwebpage

import (
	os      "os"
	strings "strings"

	util "github.com/zyxgad/go-util/util"
	ksql "github.com/zyxgad/schSvr/handles/sql"
	kses "github.com/zyxgad/schSvr/handles/sql/session"
	krbt "github.com/zyxgad/schSvr/robot"
)

var (
	sqlUserTable ksql.SqlTable
	sqlSesTable ksql.SqlTable
	sqlQuestTable ksql.SqlTable
	sqlMatchTable ksql.SqlTable
	sqlMatchUserTable ksql.SqlTable
	sqlMatchAnswerTable ksql.SqlTable
)

type SqlUserType struct{
	Id           uint32
	Username     string
	ShaPwd       string
	Password     string
	Verified     bool
	Frozen       bool
	Manager      uint32
}

func hashUserPwd(user *SqlUserType)(hashc string){
	return util.BytesToHex(util.BytesToSha256(([]byte)(user.Username + ";" + user.Password)))
}

func getUser(wheremap ksql.WhereMap)(user *SqlUserType){
	lines, err := sqlUserTable.SqlSearch(ksql.TypeMap{
		"id": ksql.TYPE_Uint32,
		"username": ksql.TYPE_String,
		"password": ksql.TYPE_String,
		"verified": ksql.TYPE_Bool,
		"frozen": ksql.TYPE_Bool,
		"manager": ksql.TYPE_Uint32,
	}, wheremap, 1)
	if err != nil || len(lines) != 1 {
		return nil
	}
	line := lines[0]
	return &SqlUserType{
		Id: util.JsonToUint32(line["id"]),
		Username: util.JsonToString(line["username"]),
		ShaPwd: util.JsonToString(line["password"]),
		Password: "",
		Verified: util.JsonToBool(line["verified"]),
		Frozen: util.JsonToBool(line["frozen"]),
		Manager: util.JsonToUint32(line["manager"]),
	}
}

func getUserById(id uint32)(user *SqlUserType){
	return getUser(ksql.WhereMap{{"id", "=", id, ""}})
}

func getUserByName(name string)(user *SqlUserType){
	return getUser(ksql.WhereMap{{"username", "=", name, ""}})
}

func updateUserById(user *SqlUserType)(err error){
	user.ShaPwd = hashUserPwd(user)
	err = sqlUserTable.SqlUpdate(ksql.Map{
		"username": user.Username,
		"password": user.ShaPwd,
		"verified": user.Verified,
		"frozen": user.Frozen,
		"manager": user.Manager,
	}, ksql.WhereMap{{"id", "=", user.Id, ""}})
	return err
}

func setUserPwd(user *SqlUserType)(err error){
	user.ShaPwd = hashUserPwd(user)
	err = sqlUserTable.SqlUpdate(ksql.Map{"password": user.ShaPwd}, ksql.WhereMap{{"id", "=", user.Id, ""}})
	return err
}

func createUser(user *SqlUserType)(err error){
	user.ShaPwd = hashUserPwd(user)
	user.Verified = false
	user.Frozen = false
	err = sqlUserTable.SqlInsert(ksql.Map{
		"id": user.Id,
		"username": user.Username,
		"password": user.ShaPwd,
		"verified": user.Verified,
		"frozen": user.Frozen,
		// "manager": user.Manager,
	})
	return err
}

func verifiedUser(user *SqlUserType)(ok bool){
	if user.Verified {
		return true
	}
	user.Verified = true
	err := sqlUserTable.SqlUpdate(ksql.Map{"verified": true}, ksql.WhereMap{{"id", "=", user.Id, ""}})
	return err == nil
}

func frozenUser(user *SqlUserType, frozen bool)(ok bool){
	user.Frozen = frozen
	err := sqlUserTable.SqlUpdate(ksql.Map{"frozen": frozen}, ksql.WhereMap{{"id", "=", user.Id, ""}})
	return err != nil
}

func setUserB64Head(b64data string, user *SqlUserType)(ok bool){
	imgpath := util.JoinPath(USER_DATA_PATH, util.JoinObjStr(user.Id), "head.png")
	if len(b64data) == 0 {
		util.RemoveFile(imgpath)
		return true
	}
	typeind := strings.Index(b64data, "data:image/")
	dataind := strings.Index(b64data, "base64,")
	if typeind == -1 || dataind == -1 {
		return false
	}
	var (
		imgbytes []byte
		err error
		file *os.File
	)
	// typestr := b64data[typeind:dataind]
	// if typestr != "png" {
	// 	return false
	// }
	b64data = b64data[dataind + 7:]
	imgbytes, err = util.DecodeBase64(b64data)
	if err != nil {
		logger.Debugln("decode user head error:", err)
		return false
	}
	file, err = os.OpenFile(imgpath, os.O_RDWR | os.O_CREATE, os.ModePerm)
	if err != nil {
		logger.Debugln("open user head file error:", err)
		return false
	}
	defer file.Close()
	_, err = file.Write(imgbytes)
	if err != nil {
		logger.Debugln("write user head error:", err)
		return false
	}
	return true
}


type sqlSource int

func (sqlSource)Init(){
	sqlSesTable = ksql.NewAutoCleanTable(sqloper.NewTable("SESSION_LINE"), 60 * 60 * 24 * 30)
	kses.SetTable(sqlSesTable)
	sqlUserTable = sqloper.NewTable("USERS")
	// logger.Debugln()
	krbt.InitCaptchaSqlTable(ksql.NewAutoCleanTable(sqloper.NewTable("CAPTCHA_DATAS"), 60 * 10))
	sqlQuestTable = sqloper.NewTable("QUESTIONS")
	sqlMatchTable = sqloper.NewTable("MATCHS")
	sqlMatchUserTable = sqloper.NewTable("MATCH_USERS")
	sqlMatchAnswerTable = sqloper.NewTable("MATCH_ANSWERS")
}
