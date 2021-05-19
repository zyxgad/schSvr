
package kpnmlogger

func Init(){
	(configSource)(0).Init()
	(loggerSource)(0).Init()
	(serverSource)(0).Init()
}

func Run(){
	StartServer()
}
