
package main

import (
	runtime "runtime"

	pkg "github.com/zyxgad/schSvr/servers/web/pkg"
)


func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){
	pkg.Init()
	pkg.Run()
}
