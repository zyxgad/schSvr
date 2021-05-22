
package kpnmproxy

import (
	http     "net/http"
	httputil "net/http/httputil"

	gin  "github.com/gin-gonic/gin"
	util "github.com/zyxgad/go-util/util"
)


type proxyResWriter struct{
	baseroute string

	writer http.ResponseWriter

	wrote bool
	wroteHeader bool
	movedres bool
}

func newProxyWriter(writer http.ResponseWriter, baseroute string)(http.ResponseWriter){
	return &proxyResWriter{
		baseroute: baseroute,
		writer: writer,
		wrote: false,
		wroteHeader: false,
		movedres: false,
	}
}

func (w *proxyResWriter)Header()(http.Header){
	return w.writer.Header()
}

func (w *proxyResWriter)Write(data []byte)(n int, err error){
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if !w.wrote {
		if w.movedres && len(data) > 0 && data[0] == '/' {
			w.writer.Write(([]byte)(w.baseroute))
		}
	}
	return w.writer.Write(data)
}

func (w *proxyResWriter)WriteHeader(statusCode int){
	if !w.wroteHeader {
		w.writer.WriteHeader(statusCode)
		w.wroteHeader = true
		switch statusCode {
		case
			http.StatusMovedPermanently,
			http.StatusFound,
			http.StatusSeeOther,
			http.StatusTemporaryRedirect,
			http.StatusPermanentRedirect:
			w.movedres = true
		}
	}
}

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
		pwriter := newProxyWriter(cont.Writer, baseroute)
		proxy.ServeHTTP(pwriter, cont.Request)
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
