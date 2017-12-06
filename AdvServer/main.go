package main

import (
	"Adventure/AdvServer/config"
	"Adventure/AdvServer/log"
	network "Adventure/AdvServer/network2"
	"Adventure/AdvServer/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var logger = log.GetLogger()

func SignalProc() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT /*, syscall.SIGUSR1, syscall.SIGUSR2*/)

	for {
		<-ch
		os.Exit(0)
		return
	}
}

func main() {
	// Init config
	if err := config.Init(); err != nil {
		fmt.Println("config init error")
		os.Exit(-1)
	}
	conf := config.GetConfig()

	// init log
	if err := log.Init(conf.LogPath, conf.LogLevel); err != nil {
		logger.Error("InitLogger failed (%v)", err)
		os.Exit(-1)
	}

	// Init service
	if err := service.Init(); err != nil {
		return
	}

	// Listen to system signal
	go SignalProc()

	// Accept connection
	listener := network.NewTCPServer()
	listener.Run()
}
