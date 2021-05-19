
package kpnmweb

import (
	ioutil   "io/ioutil"
	template "html/template"

	util  "github.com/zyxgad/go-util/util"
)


type templatesSource int

func (templatesSource)Init(){
	engine.SetFuncMap(template.FuncMap{
		"odd": func(num int)(bool){ return num % 2 == 0 },
	})

	{
		var templateFiles []string = make([]string, 0)
		basePath := util.JoinPath(RES_PATH, "html")
		var findFunc func(path string)
		findFunc = func(path string){
			finfos, err := ioutil.ReadDir(path)
			if err != nil {
				panic(err)
			}
			for _, info := range finfos {
				fpath := util.JoinPath(path, info.Name())
				if info.IsDir() {
					findFunc(fpath)
				}else{
					templateFiles = append(templateFiles, fpath)
				}
			}
		}
		findFunc(basePath)

		if len(templateFiles) > 0 {
			engine.LoadHTMLFiles(templateFiles...)
		}
	}
}
