package cli

import (
	"flag"
)

type CLIArgs struct {
	ConfigPath string
	Version    bool
}

func Arguments() CLIArgs {
	config_path := flag.String("c", "/etc/ip-monitor/config.json", "Path to configuration file")
	version := flag.Bool("v", false, "Print version and exit")
	flag.Parse()

	args := CLIArgs{ConfigPath: *config_path, Version: *version}

	return args
}
