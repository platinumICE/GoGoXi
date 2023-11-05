package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

var client *http.Client

func init_http(conf ToolConfiguration) {

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 0,
			MaxConnsPerHost:       0,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   0,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				// should not be changed
			},
		},
		Timeout: 0,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
			// this prevents Go from following redirects
		},
	}
}
