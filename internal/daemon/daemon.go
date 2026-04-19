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
	Exit      = make(chan bool, 1)
	state     = make(map[int]net.IP)
	stateLock = sync.Mutex{}
)

func Start(config cfg.Configuration) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			Run(config)
			for i := 1; i <= config.Interval; i++ {
				select {
				case <-Exit:
					stateLock.Lock()
					state = make(map[int]net.IP)
					stateLock.Unlock()
					return
				default:
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	return done
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
		ip, err = http_resolver.Resolve(config.HttpIpProvider, monitor.IPVersion, config.Timeout)
		if err != nil {
			return
		}
		logger.Debugf("Received IP address for %s from HTTP provider: %s", monitor.Name, ip)
	}
	if monitor.Type == "dns" {
		ip, err = dns_resolver.Resolve(monitor.Domain, monitor.IPVersion, config.Timeout, config.MaxCnameLookups)
		if err != nil {
			return
		}
		logger.Debugf("Received IP address for %s from DNS provider: %s", monitor.Domain, ip)
	}

	stateLock.Lock()
	oldIP, ok := state[i]
	if !ok {
		state[i] = ip
		stateLock.Unlock()
		logger.Debugf("Setting initial IP address for %s to: %s", monitor.Name, ip)
		return
	}

	if oldIP.String() != ip.String() {
		state[i] = ip
		stateLock.Unlock()
		logger.Infof("Detected IP change for %s: %s => %s", monitor.Name, oldIP, ip)
		if len(monitor.Triggers) == 0 {
			logger.Warnf("No trigger defined for %s. Doing nothing.", monitor.Name)
		} else {
			trigger.Execute(monitor.Triggers, oldIP.String(), ip.String())
		}
		return
	}

	stateLock.Unlock()
	logger.Debugf("No IP address change detected for %s: %s == %s", monitor.Name, oldIP, ip)
}
