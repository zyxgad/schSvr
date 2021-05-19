
package kpnmweb

import (
	// http     "net/http"
	time  "time"

	gin  "github.com/gin-gonic/gin"
	util "github.com/zyxgad/go-util/util"
)

func ginLogFunc()(gin.HandlerFunc){
	return func(cont *gin.Context){
		startTime := time.Now()
		cont.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := cont.Request.Method
		reqUri := cont.Request.RequestURI
		statusCode := cont.Writer.Status()
		clientIP := cont.ClientIP()
		logger.Infof("[%s %s %3d] PATH='%s' TIME=%v",
			clientIP,
			reqMethod,
			statusCode,
			reqUri,
			latencyTime,
		)
	}
}

func ginRecoverFunc(callback func(cont *gin.Context, err interface{}))(gin.HandlerFunc){
	return func(cont *gin.Context){
		defer util.RecoverErr(func(err interface{}){
			callback(cont, err)
		})
		cont.Next()
	}
}
