package main

import (
	logger "github.com/xlog4go"
	"fmt"
)

// channel of main process to quit.
var mainProcessQuit chan int

func main() {
	print("hi gopher")
	// main proc setup env
	MainProcSetup()

	mainProcessQuit = make(chan int)


	// register signal proc
	go signal_proc()

	// init log
	if err := logger.SetupLogWithConf(logFile); err != nil {
		fmt.Println("log init fail: %s", err.Error())
		return
	}
	defer logger.Close()

	// wait signal to exit
	value := <-mainProcessQuit

	logger.Info("QUIT-SIGNAL=%v", value)
	fmt.Println("MAIN-PROCESS Stopping...")
}
