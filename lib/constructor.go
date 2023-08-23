package lib

import (
	"github.com/dwnGnL/somPayment"
	"net/http"
	"time"
)

const (
	idleConnTimeoutSec = 60
	requestTimeoutSec  = 10
)

type SOM struct {
	client *http.Client
	config *somPayment.Config
}

func New(config *somPayment.Config) *SOM {
	httpClient := http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Second * idleConnTimeoutSec,
		},
		Timeout: time.Second * requestTimeoutSec,
	}

	return &SOM{
		client: &httpClient,
		config: config,
	}
}
