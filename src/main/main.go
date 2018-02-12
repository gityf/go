package main

import (
	logger "github.com/xlog4go"
	"fmt"
)

// channel of main process to quit.
var mainProcessQuit chan int

func main() {
	print("hi gopher")
	mainProcessQuit = make(chan int)


	// register signal proc
	go signal_proc()


	// wait signal to exit
	value := <-mainProcessQuit

	logger.Info("QUIT-SIGNAL=%v", value)
	fmt.Println("MAIN-PROCESS Stopping...")
}
