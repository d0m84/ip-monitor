package cfg

import (
	"encoding/json"
	"io"
	"os"

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

	for i := range config.Monitors {
		if config.Monitors[i].Domain == "" {
			config.Monitors[i].Type = "http"
		} else {
			config.Monitors[i].Type = "dns"
		}
		logger.Infof("Initialized monitor %s of type %s", config.Monitors[i].Name, config.Monitors[i].Type)
	}

	return config
}
