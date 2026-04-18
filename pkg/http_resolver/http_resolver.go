package http_resolver

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

var (
	timeout int = 10
)

func Resolve(provider string, ip_version string) (net.IP, error) {
	var dialer net.Dialer
	var client = &http.Client{Timeout: time.Second * time.Duration(timeout)}
	var tcp_version string = "tcp"

	if provider == "" {
		logger.Errorln("HTTP provider is empty")
		return nil, errors.New("http provider is empty")
	}

	switch ip_version {
	case "ip4":
		tcp_version = "tcp4"
	case "ip6":
		tcp_version = "tcp6"
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, tcp_version, addr)
	}
	client.Transport = transport

	resp, err := client.Get(provider)
	if err != nil {
		logger.Errorf("Error connecting to HTTP IP provider: %s", err)
		return nil, errors.New("http error")
	}
	defer resp.Body.Close()

	status_ok := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !status_ok {
		logger.Errorf("HTTP status error from IP provider: %s", provider)
		return nil, errors.New("status error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading HTTP IP provider response: %s", err)
		return nil, errors.New("response error")
	}

	ipString := strings.TrimSpace(string(body))
	if ipString == "" {
		logger.Errorf("Received empty HTTP response body from IP provider: %s", provider)
		return nil, errors.New("body error")
	}

	ip := net.ParseIP(ipString)

	if ip == nil {
		logger.Errorf("Received HTTP body can not be parsed as IP address: %s", ipString)
		return nil, errors.New("body error")
	}

	// Verify IP version matches requested version
	if ip_version == "ip4" && ip.To4() == nil {
		logger.Errorf("Expected IPv4 but received IPv6: %s", ip.String())
		return nil, errors.New("ip version mismatch")
	}
	if ip_version == "ip6" && ip.To4() != nil {
		logger.Errorf("Expected IPv6 but received IPv4: %s", ip.String())
		return nil, errors.New("ip version mismatch")
	}

	return ip, nil
}
