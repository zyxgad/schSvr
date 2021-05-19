
package kpnmwebpage

import (
	http  "net/http"

	gin   "github.com/gin-gonic/gin"
)


type indexPageSrc int

func (indexPageSrc)IndexPage(cont *gin.Context){
	cont.HTML(http.StatusOK, "/index.html", gin.H{
	})
}

func (page indexPageSrc)Init(){
	engine.GET("/", page.IndexPage)
}

