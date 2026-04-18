package cfg

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

type Configuration struct {
	LogLevel       string `json:"log_level"`
	LogTimestamps  bool   `json:"log_timestamps"`
	Interval       int    `json:"interval"`
	HttpIpProvider string `json:"http_ip_provider"`
	Monitors       []struct {
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

	if config.LogLevel == "debug" {
		logger.SetLevelDebug()
	} else {
		logger.SetLevelInfo()
	}

	if config.LogTimestamps {
		logger.Formatter.DisableTimestamp = false
	}

	if config.HttpIpProvider == "" {
		config.HttpIpProvider = "https://api.ipify.org"
	}

	if config.Interval <= 0 {
		logger.Warnf("Invalid interval specified: %d. Defaulting to 60 seconds", config.Interval)
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

		config.Monitors[i].Type = monitorType
		logger.Infof("Initialized monitor %s of type %s", config.Monitors[i].Name, config.Monitors[i].Type)
	}

	return config
}
