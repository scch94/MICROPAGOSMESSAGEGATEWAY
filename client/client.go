package client

import (
	"net/http"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
)

var Client http.Client

func InitHttpClient() {
	tr := &http.Transport{
		MaxIdleConns:        config.Config.Client.MaxIdleConns,
		MaxConnsPerHost:     config.Config.Client.MaxConnsPerHost,
		MaxIdleConnsPerHost: config.Config.Client.MaxConnsPerHost,
		IdleConnTimeout:     time.Duration(config.Config.Client.IdleConnTimeoutSeconds) * time.Second,
		DisableCompression:  config.Config.Client.DisableCompression,
	}
	Client = http.Client{
		Transport: tr,
		Timeout:   time.Duration(config.Config.Client.PetitionsTimeOut) * time.Second,
	}
}
