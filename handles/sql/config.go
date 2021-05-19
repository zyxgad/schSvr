
package kpnmsqlcli

import (
	bufio "bufio"
	os    "os"

	util "github.com/zyxgad/go-util/util"
)


var (
	SQL_HOST string
	SQL_PORT string
)

func init(){
	{
		var fd *os.File
		var err error
		fd, err = os.Open(util.JoinPath("/", "var", "server", ".config", "servers", "sql.txt"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		breader := bufio.NewReaderSize(fd, 1024)

		SQL_HOST = util.MustReadLine(breader)
		SQL_PORT = util.MustReadLine(breader)
	}
}
