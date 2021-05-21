
package kpnmwebpage

import (
	http "net/http"

	gin  "github.com/gin-gonic/gin"
	// util "github.com/zyxgad/go-util/util"
	// kses "github.com/zyxgad/schSvr/handles/sql/session"
	// ksql "github.com/zyxgad/schSvr/handles/sql"
)


type managerPageSrc int

func (managerPageSrc)indexGetPage(cont *gin.Context){
	cont.HTML(http.StatusOK, "manager/index.html", gin.H{
	})
}

func (page managerPageSrc)Init(){
	managerGroup := engine.Group("manager");{
		managerGroup.GET("/", page.indexGetPage)
	}
}

