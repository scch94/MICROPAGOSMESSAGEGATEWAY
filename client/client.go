package client

import (
	"net/http"
	"time"
)

var transport = &http.Transport{

	MaxIdleConns:          100,
	MaxIdleConnsPerHost:   100,
	MaxConnsPerHost:       100,
	IdleConnTimeout:       90 * time.Second, // Aumentar tiempo de espera para conexiones inactivas
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 30 * time.Second, // Tiempo de espera para los encabezados de respuesta
	DisableKeepAlives:     false,            // Temporalmente deshabilitado
}

var client = &http.Client{
	Transport: transport,
	Timeout:   30 * time.Second, // Tiempo de espera general para las solicitudes
}
