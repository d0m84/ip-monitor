package main

import (
	"github.com/d0m84/ip-monitor/internal/cfg"
	"github.com/d0m84/ip-monitor/internal/cli"
	"github.com/d0m84/ip-monitor/internal/daemon"
	"github.com/d0m84/ip-monitor/pkg/logger"
)

func main() {
	logger.Infoln("Starting ip-monitor daemon")
	cli_args := cli.Arguments()
	config := cfg.LoadConfiguration(cli_args.ConfigPath)

	if config.LogLevel != "debug" {
		logger.SetLevelInfo()
	}

	daemon.Run(config)
}
