
package kpnmproxy

import (
	context "context"
	http    "net/http"
	os      "os"
	signal  "os/signal"
	syscall "syscall"
	time    "time"

	gin     "github.com/gin-gonic/gin"
	util    "github.com/zyxgad/go-util/util"
)


var (
	engine *gin.Engine = gin.New()
	server *http.Server
)


type serverSource int

func (serverSource)Init(){
	logger.Infoln("Initing proxy")

	engine.Use(func(cont *gin.Context){
		defer util.RecoverErr(func(err interface{}){
			cont.JSON(http.StatusOK, gin.H{
				"error": "HTTP/500",
				"errorMessage": util.JoinObjStr(err),
			})
			cont.Abort()
		})
		cont.Next()
	})
	engine.NoRoute(func(cont *gin.Context){
		cont.JSON(http.StatusNotFound, gin.H{
			"error": "HTTP/404",
			"errorMessage": util.JoinObjStr("No route:", cont.Request.URL.String()),
		})
		cont.Abort()
	})

	(proxySrc)(0).Init()
}

func StartServer(){
	logger.Infoln("Proxy starting...")

	server = &http.Server{
		Addr:         HOST + ":" + PORT,
		ReadTimeout:  16 * time.Second,
		WriteTimeout: 16 * time.Second,
		Handler:      engine,
	}
	go func(){
		logger.Infof("Server run at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalln("Listening error:", err)
			panic(err)
		}
	}()

	bgcont := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		timeoutCtx, _ := context.WithTimeout(bgcont, 16 * time.Second)
		logger.Warnln("Notify sigs")
		server.Shutdown(timeoutCtx)
		logger.Warnln("Proxy shutdown")
	}
}
