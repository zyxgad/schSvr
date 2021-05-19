
package kpnmproxy

import (
	bufio "bufio"
	os    "os"

	util  "github.com/zyxgad/go-util/util"
	klog  "github.com/zyxgad/schSvr/handles/logger"
)

var (
	logger *klog.Logger
)

var (
	HOST string// = "127.0.0.1"
	PORT string// = "30000"
)

type configSource int

func readIpConfigFile(file string)(host string, port string){
	var fd *os.File
	var err error
	fd, err = os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	breader := bufio.NewReaderSize(fd, 1024)

	host = util.MustReadLine(breader)
	port = util.MustReadLine(breader)
	return
}

func (configSource)Init(){
	/*{ // read config file
		var fd *os.File
		var err error
		fd, err = os.Open(util.GetAbsPath(util.JoinPath(".config", "servers", "proxy.txt")))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		breader := bufio.NewReaderSize(fd, 1024)

		HOST = util.MustReadLine(breader)
		PORT = util.MustReadLine(breader)
	}*/
	HOST, PORT = readIpConfigFile(util.JoinPath("/", "var", "server", ".config", "servers", "proxy.txt"))
}

func init(){
	logger = klog.NewLogger("webproxy")
	util.AssertError(logger.Conn(), "Connect logger server error: ")
}
