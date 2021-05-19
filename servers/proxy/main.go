
package main

import (
	runtime "runtime"

	pkg "github.com/zyxgad/schSvr/servers/proxy/pkg"
)


func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){
	pkg.Init()
	pkg.Run()
}
