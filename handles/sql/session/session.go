
package kpnmsqlsession

import (
	time    "time"

	ksql "github.com/zyxgad/schSvr/handles/sql"
	util "github.com/zyxgad/go-util/util"
)

var (
	sqlTable ksql.SqlTable
)

func SetTable(tb ksql.SqlTable){
	sqlTable = tb
}

type SqlSessionValue struct{
	Uuid     string
	Key      string
	Value    string
	Overtime time.Time
}
type SessionMap map[string]*SqlSessionValue

func GetSessionMap(uuid string)(session SessionMap){
	session = make(SessionMap)

	lines, err := sqlTable.SqlSearch(ksql.TypeMap{
		"key": ksql.TYPE_String,
		"value": ksql.TYPE_String,
		"overtime": ksql.TYPE_String,
	}, ksql.WhereMap{{"uuid", "=", uuid, ""}}, 0)
	if err != nil {
		return nil
	}
	for _, l := range lines {
		ses := &SqlSessionValue{
			Uuid: uuid,
			Key: util.JsonToString(l["key"]),
			Value: util.JsonToString(l["value"]),
			Overtime: ksql.ParseDatetime(util.JsonToString(l["overtime"])),
		}
		session[ses.Key] = ses
	}

	return session
}

func GetSession(uuid string, key string)(value *SqlSessionValue){
	lines, err := sqlTable.SqlSearch(ksql.TypeMap{
		"value": ksql.TYPE_String,
		"overtime": ksql.TYPE_String,
	}, ksql.WhereMap{{"uuid", "=", uuid, "AND"}, {"key", "=", key, ""}}, 1)
	if err != nil || len(lines) == 0 {
		return nil
	}
	line := lines[0]

	value = &SqlSessionValue{
		Uuid: uuid,
		Key: key,
		Value: util.JsonToString(line["value"]),
		Overtime: ksql.ParseDatetime(util.JsonToString(line["overtime"])),
	}

	return value
}

func SetSessionMap(session SessionMap)(err error){
	if session == nil {
		return nil
	}

	for _, ses := range session {
		err = SetSession(ses)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetSession(value *SqlSessionValue)(err error){
	if value == nil {
		return nil
	}

	otime := value.Overtime
	if otime.Before(util.GetTimeNow()) {
		return nil
	}

	if sqlTable.HasData(ksql.WhereMap{{"uuid", "=", value.Uuid, "AND"}, {"key", "=", value.Key, ""}}) {
		err = sqlTable.SqlUpdate(ksql.Map{
			"value": value.Value,
			"overtime": ksql.FormatDatetime(otime),
		}, ksql.WhereMap{{"uuid", "=", value.Uuid, "AND"}, {"key", "=", value.Key, ""}})
	}else{
		err = sqlTable.SqlInsert(ksql.Map{
			"uuid": value.Uuid,
			"key": value.Key,
			"value": value.Value,
			"overtime": ksql.FormatDatetime(otime),
		})
	}

	return err
}

func RemoveAllSession(uuid string)(err error){
	err = sqlTable.SqlDelete(ksql.WhereMap{{"uuid", "=", uuid, ""}})
	return err
}

func RemoveSession(uuid string, key string)(err error){
	err = sqlTable.SqlDelete(ksql.WhereMap{{"uuid", "=", uuid, "AND"}, {"key", "=", key, ""}})
	return err
}

func ChangeSessionUuid(uuid string)(nuuid string, err error){
	var ok bool
	nuuid, ok = NewSessionUuid()
	if !ok {
		return "", nil
	}
	err = sqlTable.SqlUpdate(ksql.Map{"uuid": nuuid}, ksql.WhereMap{{"uuid", "=", uuid, ""}})
	return nuuid, err
}

func NewSessionUuid()(uuid string, ok bool){
	uuid = util.UUID2Hex(util.NewUUID())
	return uuid, true
}

