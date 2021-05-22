
package kpnmsql

import(
	sql   "database/sql"
	time  "time"

	util  "github.com/zyxgad/go-util/util"
)

var (
	sqldb *sql.DB = nil
)

func IsConn()(bool){
	return sqldb != nil
}

func Connect(){
	if IsConn() {
		return
	}

	var err error

	logger.Infoln("Connecting mysql...")

	// logger.Debugln("Sql host:", _SQL_HOST)
	// logger.Debugln("Sql port:", _SQL_PORT)
	// logger.Debugln("Sql user:", _SQL_USER)
	// logger.Debugf ("Sql pwd(sha): **%s", util.BytesToHex(util.BytesToSha256(([]byte)(_SQL_PWD))[2:]) )
	// logger.Debugln("Sql database:", _SQL_BASE)
	// logger.Debugln("Sql charset:", _SQL_CSET)

	dbDSN := util.FormatStr(
		"%s:%s@tcp(%s:%s)/%s?charset=%s", //"&parseTime=true",
		_SQL_USER, _SQL_PWD, _SQL_HOST, _SQL_PORT, _SQL_BASE, _SQL_CSET)

	sqldb, err = sql.Open("mysql", dbDSN)
	if err != nil {
		logger.Fatalln("Connect MySql error:", err)
		panic(err)
		return
	}
	sqldb.SetMaxOpenConns(64)
	sqldb.SetConnMaxLifetime(120 * time.Second)

	err = sqldb.Ping()
	if err != nil {
		logger.Fatalln("Ping MySql error:", err)
		panic(err)
		return
	}
	logger.Warnln("Connect mysql succeed")
}

func CloseConn()(error){
	if sqldb != nil {
		err := sqldb.Close()
		return err
	}
	return nil
}


type sqlconnSource int

func (sqlconnSource)Init(){
	Connect()
}

