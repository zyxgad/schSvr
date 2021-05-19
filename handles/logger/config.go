
package kpnmlogcli

import (
	bufio "bufio"
	os    "os"

	util "github.com/zyxgad/go-util/util"
)


var (
	LOGGER_HOST string
	LOGGER_PORT string
)

func init(){
	{
		var fd *os.File
		var err error
		fd, err = os.Open(util.JoinPath("/", "var", "server", ".config", "servers", "logger.txt"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		breader := bufio.NewReaderSize(fd, 1024)

		LOGGER_HOST = util.MustReadLine(breader)
		LOGGER_PORT = util.MustReadLine(breader)
	}
}
