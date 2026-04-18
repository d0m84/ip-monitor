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

var (
	timeout int = 10
)

func CheckIfCNAME(domain string) (string, bool, error) {
	target, err := net.LookupCNAME(domain)
	if err != nil {
		return "", true, err
	} else if target == domain {
		logger.Debugf("Record is not a CNAME: %s == %s", domain, target)
		return target, false, nil
	} else {
		logger.Debugf("Record is a CNAME: %s != %s", domain, target)
		return target, true, nil
	}
}

func FindFinalTarget(domain string) (string, error) {
	var err error
	var target string = domain
	var is_cname bool
	for i := 0; i < 2; i++ {
		target, is_cname, err = CheckIfCNAME(target)
		if err != nil {
			logger.Errorf("Error checking if %s is a CNAME: %s", domain, err)
			return "", err
		}
		if !is_cname {
			logger.Debugf("Final target for domain %s is %s", domain, target)
			return target, nil
		}
	}
	logger.Errorf("Maximum CNAME lookup limit reached for %s", domain)
	return "", errors.New("dns cname lookup limit reached")
}

func FindNameServers(domain string) ([]*net.NS, error) {
	domain = strings.TrimSuffix(domain, ".")
	domain_parts := strings.Split(domain, ".")
	for i := range domain_parts {
		t := domain_parts[i:len(domain_parts):len(domain_parts)]
		d := strings.Join(t, ".")

		nameservers, err := net.LookupNS(d)
		if err == nil {
			return nameservers, nil
		}
	}
	return nil, errors.New("dns resolve authoritative error")
}

func LookupAuthorative(domain string, ip_version string) ([]net.IP, error) {
	nameservers, err := FindNameServers(domain)
	if err != nil {
		logger.Errorf("Unable to detect authoritative nameservers for %s", domain)
		return nil, errors.New("dns resolve authoritative error")
	}

	if len(nameservers) == 0 {
		logger.Errorf("No nameservers found for %s", domain)
		return nil, errors.New("dns resolve authoritative error")
	}

	// Resolve authoritative nameserver IP with timeout
	nsCtx, nsCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer nsCancel()

	ns_ips, err := net.DefaultResolver.LookupIPAddr(nsCtx, nameservers[rand.Intn(len(nameservers))].Host)
	if err != nil {
		logger.Errorf("Error resolving IP addresses of authoritative DNS server for %s: %s", domain, err)
		return nil, errors.New("dns lookup authoritative error")
	}

	if len(ns_ips) == 0 {
		logger.Errorf("No IP addresses found for authoritative DNS server of %s", domain)
		return nil, errors.New("dns lookup authoritative error")
	}

	ns_ip := ns_ips[0].IP
	var nameserver string
	if ns_ip.To4() != nil {
		nameserver = fmt.Sprintf("%s:53", ns_ip)
	} else {
		nameserver = fmt.Sprintf("[%s]:53", ns_ip)
	}
	logger.Debugf("Using nameserver %s for %s", nameserver, domain)

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * time.Duration(timeout),
			}
			return d.DialContext(ctx, network, nameserver)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	ips, err := r.LookupIP(ctx, ip_version, domain)
	if err != nil {
		logger.Errorf("Error resolving domain %s: %s", domain, err)
		return nil, errors.New("dns hosts error")
	}

	// Verify all IPs match requested version
	for _, ip := range ips {
		if ip_version == "ip4" && ip.To4() == nil {
			logger.Errorf("Expected IPv4 but received IPv6 for %s: %s", domain, ip.String())
			return nil, errors.New("ip version mismatch")
		}
		if ip_version == "ip6" && ip.To4() != nil {
			logger.Errorf("Expected IPv6 but received IPv4 for %s: %s", domain, ip.String())
			return nil, errors.New("ip version mismatch")
		}
	}

	return ips, nil
}

func Resolve(domain string, ip_version string) (net.IP, error) {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		logger.Errorln("DNS domain is empty")
		return nil, errors.New("dns domain empty")
	}

	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	target, err := FindFinalTarget(domain)
	if err != nil {
		return nil, errors.New("dns cname lookup error")
	}

	ips, err := LookupAuthorative(target, ip_version)
	if err != nil {
		return nil, errors.New("dns authoritative error")
	}
	logger.Debugf("Resolved IP addresses for %s: %s", domain, ips)

	if len(ips) == 0 {
		logger.Errorf("No host records found for %s", domain)
		return nil, errors.New("dns hosts error")
	}

	if len(ips) > 1 {
		logger.Errorf("Received multiple host entries for %s: %s", domain, ips)
		return nil, errors.New("multiple host records found")
	}

	return ips[0], nil
}
