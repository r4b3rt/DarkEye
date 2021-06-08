package main

import (
	"github.com/zsdevX/DarkEye/common"
	"runtime"
)

func main() {
	initializer()
	superScanRuntimeOptions.Start()

}

func initializer() {
	//初始化系统
	common.SetRLimit()
	runtime.GOMAXPROCS(runtime.NumCPU())
	superScanRuntimeOptions.Init()
	//  debug/pprof
	/*
		go func() {
			fmt.Println(http.ListenAndServe("localhost:10000", nil))
		}()
	*/
}


