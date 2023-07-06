package daemon

import (
	"net"
	"sync"
	"time"

	"github.com/d0m84/ip-monitor/internal/cfg"
	"github.com/d0m84/ip-monitor/internal/trigger"
	"github.com/d0m84/ip-monitor/pkg/dns_resolver"
	"github.com/d0m84/ip-monitor/pkg/http_resolver"
	"github.com/d0m84/ip-monitor/pkg/logger"
)

var (
	Exit  = make(chan bool)
	state = make(map[int]net.IP)
)

func Start(config cfg.Configuration) {
	go func() {
		for {
			Run(config)
			for i := 1; i <= config.Interval; i++ {
				select {
				case <-Exit:
					state = make(map[int]net.IP)
					return
				default:
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
}

func Run(config cfg.Configuration) {
	var wg = &sync.WaitGroup{}
	for i := range config.Monitors {
		wg.Add(1)
		go func(i int) {
			Monitor(&config, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func Monitor(config *cfg.Configuration, i int) {
	var ip net.IP
	var err error
	monitor := config.Monitors[i]
	if monitor.Type == "http" {
		ip, err = http_resolver.Resolve(config.HttpIpProvider, monitor.IPVersion)
		if err != nil {
			return
		}
		logger.Debugf("Received IP address for %s from HTTP provider: %s", monitor.Name, ip)
	}
	if monitor.Type == "dns" {
		ip, err = dns_resolver.Resolve(monitor.Domain, monitor.IPVersion)
		if err != nil {
			return
		}
		logger.Debugf("Received IP address for %s from DNS provider: %s", monitor.Domain, ip)
	}

	_, ok := state[i]
	if !ok {
		logger.Debugf("Setting initial IP address for %s to: %s", monitor.Name, ip)
		state[i] = ip
	} else if state[i].String() != ip.String() {
		logger.Infof("Detected IP change for %s: %s => %s", monitor.Name, state[i], ip)
		if len(monitor.Triggers) == 0 {
			logger.Warnf("No trigger defined for %s. Doing nothing.", monitor.Name)
		} else {
			trigger.Execute(monitor.Triggers, state[i].String(), ip.String())
		}
		state[i] = ip
	} else {
		logger.Debugf("No IP address change detected for %s: %s == %s", monitor.Name, state[i], ip)
	}
}
