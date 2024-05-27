package client

import (
	"net"
	"net/http"
	"time"
)

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,

	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:        300,
	IdleConnTimeout:     90 * time.Second,
	TLSHandshakeTimeout: 10 * time.Second,
}

var client = &http.Client{
	Transport: transport,
}
