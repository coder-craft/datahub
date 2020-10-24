package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"../core"
	"../core/common/zlog"
	"../core/conf"
)

const (
	ver = "1.0.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	zlog.Info("service start", zlog.String("Version", ver))
	core.Start()
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		zlog.Info("service get a signal", zlog.String("Signal", s.String()))
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			core.Stop()
			zlog.Info("service exit", zlog.String("Version", ver))
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
