
package kpnmwebpage

import (
	gin  "github.com/gin-gonic/gin"
	klog "github.com/zyxgad/schSvr/handles/logger"
	ksql "github.com/zyxgad/schSvr/handles/sql"
)

var (
	logger *klog.Logger
	sqloper *ksql.Operator
	engine *gin.Engine

	RES_PATH string
	USER_DATA_PATH string
)


func InitAllPage(eng *gin.Engine, logger_ *klog.Logger, sqloper_ *ksql.Operator, engine_ *gin.Engine, respath string){
	logger = logger_
	sqloper = sqloper_
	engine = engine_
	RES_PATH = respath
	USER_DATA_PATH = "/var/server/userdata"


	(sqlSource)(0).Init()

	(resPageSrc)(0).Init()
	(indexPageSrc)(0).Init()
	(userPageSrc)(0).Init()
	(questPageSrc)(0).Init()
	(managerPageSrc)(0).Init()
}
