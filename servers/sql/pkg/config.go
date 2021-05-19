
package kpnmsql

import (
	bufio "bufio"
	os    "os"

	_     "github.com/go-sql-driver/mysql"
	util   "github.com/zyxgad/go-util/util"
	klog  "github.com/zyxgad/schSvr/handles/logger"
)

var (
	logger *klog.Logger
)

var (
	HOST string// = "127.0.0.1"
	PORT string// = "30000"

	MAX_CONN uint32 = 32

	_SQL_HOST string
	_SQL_PORT string
	_SQL_USER string
	_SQL_PWD  string
	_SQL_BASE string
	_SQL_CSET string
)

type configSource int

func (configSource)Init(){
	logger.Infoln("The sql server is initializing")

	{// read sql config
		var fd *os.File
		var err error
		fd, err = os.Open(util.JoinPath("/", "var", "server", ".config", "sql_conf.txt"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		breader := bufio.NewReaderSize(fd, 1024)

		_SQL_HOST = util.MustReadLine(breader)
		_SQL_PORT = util.MustReadLine(breader)
		_SQL_USER = util.MustReadLine(breader)
		_SQL_PWD  = util.MustReadLine(breader)
		_SQL_BASE = util.MustReadLine(breader)
		_SQL_CSET = util.MustReadLine(breader)
	}
	{// read config file
		var fd *os.File
		var err error
		fd, err = os.Open(util.JoinPath("/", "var", "server", ".config", "servers", "sql.txt"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		breader := bufio.NewReaderSize(fd, 1024)

		HOST = util.MustReadLine(breader)
		PORT = util.MustReadLine(breader)
	}
}

func init(){
	logger = klog.NewLogger("svrsql")
	util.AssertError(logger.Conn(), "Connect logger server error: ")
}
