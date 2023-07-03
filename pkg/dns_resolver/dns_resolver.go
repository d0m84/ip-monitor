package dns_resolver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

func LookupLocal(domain string, ip_version string) ([]net.IP, error) {
	ips, err := net.DefaultResolver.LookupIP(context.Background(), ip_version, domain)
	if err != nil {
		logger.Errorf("Error resolving domain name: %s", err)
		return nil, errors.New("dns cache error")
	}

	return ips, nil
}

func FindNameServers(domain string) ([]*net.NS, error) {
	domainParts := strings.Split(domain, ".")

	for i := range domainParts {
		t := domainParts[i:len(domainParts):len(domainParts)]
		d := strings.Join(t, ".")

		nameservers, err := net.LookupNS(d)
		if err == nil {
			return nameservers, nil
		}
	}

	return nil, errors.New("dns resolve authorative error")
}

func LookupAuthorative(domain string, ip_version string) ([]net.IP, error) {
	nameservers, err := FindNameServers(domain)
	if err != nil {
		logger.Errorf("Unable to get authorative nameservers")
		return nil, errors.New("dns resolve authorative error")
	}

	ns_ips, err := net.LookupIP(nameservers[rand.Intn(len(nameservers))].Host)
	if err != nil {
		logger.Errorf("Error resolving IP addresses of authorative DNS server: %s", err)
		return nil, errors.New("dns lookup authorative error")
	}

	nameserver := fmt.Sprintf("%s:53", ns_ips[0])
	logger.Debugf("Using nameserver %s", nameserver)

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, nameserver)
		},
	}

	ips, err := r.LookupIP(context.Background(), ip_version, domain)
	if err != nil {
		logger.Errorf("Error resolving domain name: %s", err)
		return nil, errors.New("dns hosts error")
	}

	return ips, nil
}

func Resolve(domain string, ip_version string) (net.IP, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		logger.Errorf("Error checking for cname: %s", err)
		return nil, errors.New("dns cname error")
	}

	var ips []net.IP

	if cname == fmt.Sprintf("%s.", domain) {
		logger.Debugf("Using authorative DNS servers for %s", domain)
		ips, err = LookupAuthorative(domain, ip_version)
		if err != nil {
			return nil, errors.New("dns authorative error")
		}
	} else {
		logger.Debugf("Using local DNS servers for %s", domain)
		ips, err = LookupLocal(domain, ip_version)
		if err != nil {
			return nil, errors.New("dns cache error")
		}
	}

	logger.Debugf("Resolved IP addresses: %s", ips)

	if len(ips) > 1 {
		logger.Warnf("Received multiple host entries for domain %s. Using first entry.", domain)
	}

	return ips[0], nil
}
