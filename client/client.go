package client

import (
	"net/http"
	"time"
)

var Client http.Client

func InitHttpClient() {
	// TODO: Pasar a una config
	maxIdleConns := 10
	maxConnsPerHost := 100
	maxIdleConnsPerHost := 10
	idleConnTimeoutSeconds := 30
	disableCompression := true
	requestTimeout := 30
	tr := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxConnsPerHost:     maxConnsPerHost,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(idleConnTimeoutSeconds) * time.Second,
		DisableCompression:  disableCompression,
	}
	Client = http.Client{
		Transport: tr,
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}
}
