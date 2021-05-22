
package kpnmweb

import (
	context "context"
	http    "net/http"
	os      "os"
	signal  "os/signal"
	syscall "syscall"
	time    "time"

	gin     "github.com/gin-gonic/gin"
	util    "github.com/zyxgad/go-util/util"

	kwebpage "github.com/zyxgad/schSvr/servers/web/pkg/pages"
)

var (
	engine *gin.Engine = gin.New()
	server *http.Server
)


type serverSource int

func (serverSource)Init(){
	logger.Infoln("Initing webs")

	engine.Use(ginLogFunc(), ginRecoverFunc(func(cont *gin.Context, err interface{}){
		logger.Errorln("Error:", err)
		logger.Errorln("Error:", util.GetStack(0))
		cont.HTML(http.StatusInternalServerError, "errs/E500.html", gin.H{
			"cause": util.JoinObjStr(err),
			"stacks": util.GetStacks(0),
		})
		cont.Abort()
	}))
	engine.NoRoute(func(cont *gin.Context) {
		cont.HTML(http.StatusNotFound, "errs/E404.html", gin.H{
			"path": cont.Request.URL.String(),
		})
		cont.Abort()
	})

	kwebpage.InitAllPage(engine, logger, sqloper, engine, RES_PATH)
}

func StartServer(){
	logger.Warnln("Web server starting...")

	// logger.Infoln("Gin", gin.Version)
	// logger.Infof("Server run at %s:%s", HOST, PORT)

	server = &http.Server{
		Addr:         HOST + ":" + PORT,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      engine,
	}
	go func(){
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(util.JoinObjStr("Listening error:", err))
		}
	}()

	bgcont := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		timeoutCtx, _ := context.WithTimeout(bgcont, 10 * time.Second)
		logger.Warnln("Notify sigs")
		server.Shutdown(timeoutCtx)
		logger.Warnln("Http shutdown")
	}
}
