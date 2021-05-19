
package kpnmutil

import (
	strings "strings"
	os      "os"
)

var (
	_RUN_MODE string
)

func GetRunMode()(string){
	return _RUN_MODE
}

func IsDebug()(bool){
	return _RUN_MODE == "DEBUG"
}

func init(){
	_RUN_MODE = strings.ToUpper(os.Getenv("SCH_SVR_MODE"))
}
