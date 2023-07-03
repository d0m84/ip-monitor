package daemon

import (
	"net"
	"time"

	"github.com/d0m84/ip-monitor/internal/cfg"
	"github.com/d0m84/ip-monitor/internal/trigger"
	"github.com/d0m84/ip-monitor/pkg/dns_resolver"
	"github.com/d0m84/ip-monitor/pkg/http_resolver"
	"github.com/d0m84/ip-monitor/pkg/logger"
)

func Run(config cfg.Configuration) {
	ip_state := make(map[int]net.IP)

	for i := range config.Monitors {
		ip_state[i] = nil
	}

	for {
		for i := range config.Monitors {
			var ip net.IP
			var err error
			monitor := config.Monitors[i]
			if monitor.Type == "http" {
				ip, err = http_resolver.Resolve(config.HttpIpProvider, monitor.IPVersion)
				if err != nil {
					continue
				}
				logger.Debugf("Received IP address for %s from HTTP provider: %s", monitor.Name, ip)
			}
			if monitor.Type == "dns" {
				ip, err = dns_resolver.Resolve(monitor.Domain, monitor.IPVersion)
				if err != nil {
					continue
				}
				logger.Debugf("Received IP address for %s from DNS provider: %s", monitor.Domain, ip)
			}

			if ip_state[i] == nil {
				logger.Debugf("Setting initial IP address to: %s", ip)
				ip_state[i] = ip
			} else if ip_state[i].String() != ip.String() {
				logger.Infof("Detected IP change for %s: %s => %s", monitor.Name, ip_state[i], ip)
				if len(monitor.Triggers) == 0 {
					logger.Warnf("No trigger defined for %s. Doing nothing.", monitor.Name)
				} else {
					trigger.Execute(monitor.Triggers, ip_state[i].String(), ip.String())
				}
				ip_state[i] = ip
			} else {
				logger.Debugf("No IP address change detected: %s == %s", ip_state[i], ip)
			}
		}

		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
