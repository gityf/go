package main

import (
	"global"
	"runtime"
)

func MainProcSetup() {
	global.MainProcIsRunning = 1
	runtime.GOMAXPROCS(runtime.NumCPU())
}