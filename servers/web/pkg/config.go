
package kpnmweb

import (
	bufio "bufio"
	os    "os"

	util  "github.com/zyxgad/go-util/util"
	klog  "github.com/zyxgad/schSvr/handles/logger"
	ksql  "github.com/zyxgad/schSvr/handles/sql"
)

var (
	logger *klog.Logger
	sqloper *ksql.Operator
)

var (
	RES_PATH = util.JoinPath("/", "var", "server", "static")

	HOST string// = "127.0.0.1"
	PORT string// = "30000"

	MAX_CONN uint32 = 128
)

type configSource int

func (configSource)Init(){
	logger.Infoln("The web server is initializing")
	{ // read config file
		var fd *os.File
		var err error
		fd, err = os.Open(util.GetAbsPath(util.JoinPath("/", "var", "server", ".config", "servers", "web.txt")))
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
	logger = klog.NewLogger("svrweb")
	sqloper = ksql.NewOperator("svrweb")
}
