package main

import (
	"global"
	"runtime"
	"path"
	"os"
	"fmt"
	"config"
)

func MainProcSetup() {
	global.MainProcIsRunning = 1
	runtime.GOMAXPROCS(runtime.NumCPU())
	ConfSetup()
}

var logFile = "./conf/log.json"
var confFile = "./conf/service.json"

func ConfSetup() {
	dirPath := path.Dir(os.Args[0])
	logFile = dirPath + "/../" + logFile
	confFile = dirPath + "/../" + confFile
	fmt.Println("logFile:", logFile)
	fmt.Println("confFile:", confFile)
	var err error
	if err = config.ParseConf(confFile); err != nil {
		fmt.Println("conf init fail: %s", err.Error())
		return
	}
}
