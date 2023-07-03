package cli

import (
	"flag"
)

type CLIArgs struct {
	ConfigPath string
}

func Arguments() CLIArgs {
	config_path := flag.String("c", "/etc/ip-monitor/config.json", "Path to configuration file")
	flag.Parse()

	args := CLIArgs{ConfigPath: *config_path}

	return args
}
