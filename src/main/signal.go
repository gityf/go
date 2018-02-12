package main

import (
	"os"
	"os/signal"
	"syscall"
	logger "github.com/xlog4go"
)

func signal_proc() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGALRM, syscall.SIGTERM, syscall.SIGUSR1)

	// Block until a signal is received.
	sig := <-c

	logger.Warn("Signal received: %v", sig)

	// TODO close server gracefully.

	logger.Warn("send quit signal")
	mainProcessQuit <- 1
}
