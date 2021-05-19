
package kpnmwebpage

import (
	http    "net/http"
	os      "os"

	gin    "github.com/gin-gonic/gin"
	util  "github.com/zyxgad/go-util/util"
)


func writeHttpFile(cont *gin.Context, path string, ctype string)(error){
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	cont.Writer.Header().Set("Content-Type", ctype)
	cont.Status(http.StatusOK)
	util.MustCopyWR(cont.Writer, fd)
	return nil
}

type resPageSrc int

func (resPageSrc)redirectPage(cont *gin.Context){
	url := cont.Query("url")
	cont.HTML(http.StatusOK, "utilp/redirect.html", gin.H{
		"url": url,
	})
}

func (resPageSrc)cssGetterPage(cont *gin.Context){
	respath := cont.Param("respath")
	mypath := util.JoinPathWithoutAbs(RES_PATH, "css", respath)
	if util.IsFile(mypath){
		err := writeHttpFile(cont, mypath, "text/css")
		if err == nil {
			return
		}
	}
	if util.IsFile(mypath + ".min.css") {
		if writeHttpFile(cont, mypath + ".min.css", "text/css") == nil {
			return
		}
	}else if util.IsFile(mypath + ".css") {
		if writeHttpFile(cont, mypath + ".css", "text/css") == nil {
			return
		}
	}
	cont.String(http.StatusNotFound, "/* Can not found css file '%s' */", respath)
}

func (resPageSrc)imgGetterPage(cont *gin.Context){
	respath := cont.Param("respath")
	mypath := util.JoinPathWithoutAbs(RES_PATH, "images", respath)
	if util.IsFile(mypath){
		fd, err := os.Open(mypath)
		if err == nil {
			// cont.Writer.Header().Set("Content-Type", "")
			cont.Status(http.StatusOK)
			util.MustCopyWR(cont.Writer, fd)
			return
		}
	}
	cont.Status(http.StatusNotFound)
}

func (resPageSrc)jsGetterPage(cont *gin.Context){
	respath := cont.Param("respath")
	mypath := util.JoinPathWithoutAbs(RES_PATH, "js", respath)
	if util.IsFile(mypath){
		err := writeHttpFile(cont, mypath, "text/javascript")
		if err == nil {
			return
		}
	}
	if util.IsFile(mypath + ".min.js") {
		if writeHttpFile(cont, mypath + ".min.js", "text/javascript") == nil {
			return
		}
	}else if util.IsFile(mypath + ".js") {
		if writeHttpFile(cont, mypath + ".js", "text/javascript") == nil {
			return
		}
	}
	cont.String(http.StatusNotFound, "// Can not found js file '%s'", respath)
}

func (page resPageSrc)Init(){
	engine.GET("/redirectto", page.redirectPage)

	resGroup := engine.Group("static");{
		resGroup.GET("/css/*respath", page.cssGetterPage)
		resGroup.GET("/images/*respath", page.imgGetterPage)
		resGroup.GET("/js/*respath", page.jsGetterPage)
	}
}

