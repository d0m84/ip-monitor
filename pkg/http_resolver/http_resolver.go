package http_resolver

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/d0m84/ip-monitor/pkg/logger"
)

func Resolve(provider string, ip_version string) (net.IP, error) {

	var zeroDialer net.Dialer
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	var tcpVersion string = ip_version

	if ip_version == "ip4" {
		tcpVersion = "tcp4"
	} else if ip_version == "ip6" {
		tcpVersion = "tcp6"
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, tcpVersion, addr)
	}
	httpClient.Transport = transport

	resp, err := httpClient.Get(provider)
	if err != nil {
		logger.Errorf("Error connecting to HTTP IP provider: %s", err)
		return nil, errors.New("http error")
	}
	defer resp.Body.Close()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		logger.Errorf("HTTP status error from IP provider: %s", provider)
		return nil, errors.New("status error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading HTTP IP provider response: %s", err)
		return nil, errors.New("response error")
	}

	ip := net.ParseIP(string(body))

	if ip == nil {
		logger.Errorf("Received HTTP body can not be parsed as IP address: %s", string(body))
		return nil, errors.New("body error")
	}

	return ip, nil
}
