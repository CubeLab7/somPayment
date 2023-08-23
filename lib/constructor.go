package lib

import (
	"github.com/dwnGnL/somPayment"
	"net/http"
	"time"
)

type SOM struct {
	client *http.Client
	config *somPayment.Config
}

func New(config *somPayment.Config) *SOM {
	httpClient := http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Second * time.Duration(config.IdleConnTimeoutSec),
		},
		Timeout: time.Second * time.Duration(config.RequestTimeoutSec),
	}

	return &SOM{
		client: &httpClient,
		config: config,
	}
}
