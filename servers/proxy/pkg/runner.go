
package kpnmproxy

import (
	os    "os"
	util  "github.com/zyxgad/go-util/util"
)

func Init(){
	defer util.RecoverErr(func(err interface{}){
		logger.Fatalln(util.JoinObjStr(err))
		logger.Fatalln(util.GetStack(0))
		os.Exit(-1)
	})
	(configSource)(0).Init()
	(serverSource)(0).Init()
}

func Run(){
	defer util.RecoverErr(func(err interface{}){
		logger.Fatalln(util.JoinObjStr(err))
		logger.Fatalln(util.GetStack(0))
	})
	StartServer()
}
