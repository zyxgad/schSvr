
package kpnmproxy

import (
	http     "net/http"
	httputil "net/http/httputil"

	gin  "github.com/gin-gonic/gin"
	util "github.com/zyxgad/go-util/util"
)


func newProxy(target string, baseroute string)(func(cont *gin.Context)){
	return func(cont *gin.Context){
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.Host = target
				req.URL.Scheme = "http"
				req.URL.Host = target
				req.URL.Path = cont.Param("route")
			},
			ErrorHandler: func(res http.ResponseWriter, req *http.Request, err error){
				res.WriteHeader(http.StatusInternalServerError)
				res.Write(([]byte)(util.EncodeJson(util.JsonType{
					"error": "HTTP/500",
					"errorMessage": "proxy panic",
					"cause": err.Error(),
				})))
			},
		}
		proxy.ServeHTTP(cont.Writer, cont.Request)
	}
}


type proxySrc int

func (proxySrc)Init(){
	engine.GET("/", func(cont *gin.Context){
		cont.Redirect(http.StatusMovedPermanently, "/web")
	})
	webGroup := engine.Group("web");{
		host, port := readIpConfigFile(util.JoinPath("/", "var", "server", ".config", "servers", "web.txt"))
		webGroup.Any("/*route", newProxy(host + ":" + port, "web"))
	}
}
