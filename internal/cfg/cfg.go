package cfg

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

type Configuration struct {
	LogLevel        string `json:"log_level"`
	LogTimestamps   bool   `json:"log_timestamps"`
	Interval        int    `json:"interval"`
	Timeout         int    `json:"timeout"`
	MaxCnameLookups int    `json:"max_cname_lookups"`
	HttpIpProvider  string `json:"http_ip_provider"`
	Monitors        []struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Domain    string   `json:"domain"`
		IPVersion string   `json:"ip_version"`
		Triggers  []string `json:"triggers"`
	} `json:"monitors"`
}

func LoadConfiguration(cfgFile string) Configuration {
	file, err := os.Open(cfgFile)
	if err != nil {
		logger.Fatalln("Unable to open config:", err)
	}
	defer file.Close()

	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Fatalln("Unable to read config:", err)
	}

	var config Configuration
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		logger.Fatalln("Unable to parse config:", err)
	}

	switch strings.ToLower(config.LogLevel) {
	case "debug":
		logger.SetLevelDebug()
	case "info":
		logger.SetLevelInfo()
	case "warn":
		logger.SetLevelWarn()
	case "error":
		logger.SetLevelError()
	default:
		logger.Warnf("Invalid log level specified: %s. Defaulting to info", config.LogLevel)
		logger.SetLevelInfo()
	}

	if config.LogTimestamps {
		logger.Formatter.DisableTimestamp = false
	}

	if config.HttpIpProvider == "" {
		config.HttpIpProvider = "https://api.ipify.org"
	}

	if config.Timeout <= 0 {
		logger.Debugf("No or invalid timeout specified: %d. Defaulting to 10 seconds", config.Timeout)
		config.Timeout = 10
	}

	if config.MaxCnameLookups <= 0 {
		logger.Debugf("No or invalid max_cname_lookups specified: %d. Defaulting to 2", config.MaxCnameLookups)
		config.MaxCnameLookups = 2
	}

	if config.Interval <= 0 {
		logger.Warnf("No or invalid interval specified: %d. Defaulting to 60 seconds", config.Interval)
		config.Interval = 60
	}

	if len(config.Monitors) == 0 {
		logger.Fatalln("No monitor configured")
	}

	for i := range config.Monitors {
		monitorType := strings.ToLower(config.Monitors[i].Type)
		if monitorType == "" {
			if config.Monitors[i].Domain == "" {
				monitorType = "http"
			} else {
				monitorType = "dns"
			}
		}
		if monitorType != "http" && monitorType != "dns" {
			logger.Fatalln("Unsupported monitor type:", config.Monitors[i].Type)
		}
		if monitorType == "dns" && config.Monitors[i].Domain == "" {
			logger.Fatalln("DNS monitors must have a domain specified")
		}
		if monitorType == "http" && config.Monitors[i].Domain != "" {
			logger.Warnf("HTTP monitors should not have a domain specified. Ignoring domain for monitor %d", i)
		}
		config.Monitors[i].Type = monitorType

		ipVersion := strings.ToLower(config.Monitors[i].IPVersion)
		if ipVersion != "ip4" && ipVersion != "ip6" {
			logger.Fatalln("Unsupported IP version:", config.Monitors[i].IPVersion)
		}

		name := config.Monitors[i].Name
		if name == "" {
			if config.Monitors[i].Domain == "" {
				logger.Warnf("No name specified for monitor %d and no domain provided. Using default name: monitor%d", i, i)
				name = "monitor" + strconv.Itoa(i)
			} else {
				logger.Warnln("No name specified for monitor. Using domain:", config.Monitors[i].Domain)
				name = config.Monitors[i].Domain
			}
			config.Monitors[i].Name = name
		}

		logger.Infof("Initialized monitor %s of type %s", name, config.Monitors[i].Type)
	}

	return config
}
