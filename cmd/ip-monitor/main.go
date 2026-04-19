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

var (
	version string = "" // Set at build time using -ldflags "-X main.version=$(VERSION)"
	date    string = "" // Set at build time using -ldflags "-X main.date=$(DATE)"
)

func main() {
	cli_args := cli.Arguments()

	if cli_args.Version {
		println("IP-Monitor version:", version, "- Build date:", date)
		os.Exit(0)
	}

	config := cfg.LoadConfiguration(cli_args.ConfigPath)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	logger.Infoln("Starting IP-Monitor Daemon")
	var daemonDone <-chan struct{}
	daemonDone = daemon.Start(config)

	for s := range sigs {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			logger.Infoln("Stopping IP-Monitor Daemon")
			daemon.Exit <- true
			<-daemonDone
			os.Exit(0)
		case syscall.SIGHUP:
			logger.Infoln("Reloading IP-Monitor Daemon")
			daemon.Exit <- true
			<-daemonDone
			// Drain any residual signal from the Exit channel
			select {
			case <-daemon.Exit:
			default:
			}
			config = cfg.LoadConfiguration(cli_args.ConfigPath)
			daemonDone = daemon.Start(config)
		}
	}

}
