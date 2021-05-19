
package kpnmlogger

import (
	os    "os"
	bufio "bufio"

	util  "github.com/zyxgad/go-util/util"
)

var (
	HOST string// = "127.0.0.1"
	PORT string// = "30000"

	MAX_CONN uint32 = 32
)


type configSource int

func (configSource)Init(){
	{
		var fd *os.File
		var err error
		fd, err = os.Open(util.JoinPath("/", "var", "server", ".config", "servers", "logger.txt"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		breader := bufio.NewReaderSize(fd, 1024)

		HOST = util.MustReadLine(breader)
		PORT = util.MustReadLine(breader)
	}
}
