
package kpnmwebpage


func (questPageSrc)matchcheckGetPage(cont *gin.Context){
	muserid := cont.Query("muserid")
	cont.HTML(http.StatusOK, "quest/matchcheck.html", gin.H{
		"muserid": muserid,
	})
}

func (page questPageSrc)Init(){
	userGroup := engine.Group("quest");{
	}
}

