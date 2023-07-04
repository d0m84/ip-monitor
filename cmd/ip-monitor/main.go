package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/d0m84/ip-monitor/internal/cfg"
	"github.com/d0m84/ip-monitor/internal/cli"
	"github.com/d0m84/ip-monitor/internal/daemon"
	"github.com/d0m84/ip-monitor/pkg/logger"
)

func main() {
	cli_args := cli.Arguments()
	config := cfg.LoadConfiguration(cli_args.ConfigPath)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	logger.Infoln("Starting IP-Monitor Daemon")
	daemon.Start(config)

	for {
		select {
		case s := <-sigs:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				logger.Infoln("Stopping IP-Monitor Daemon")
				os.Exit(0)
			case syscall.SIGHUP:
				logger.Infoln("Reloading IP-Monitor Daemon")
				daemon.Exit <- true
				config = cfg.LoadConfiguration(cli_args.ConfigPath)
				daemon.Start(config)
			}
		}
	}

}
