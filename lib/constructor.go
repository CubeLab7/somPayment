package lib

import (
	"context"
	"github.com/dwnGnL/somPayment/config"
	"net/http"
	"time"
)

const (
	idleConnTimeoutSec = 60
	requestTimeoutSec  = 10
)

type SOM struct {
	client *http.Client
	config *config.Config
}

func New(ctx *context.Context, config *config.Config) *SOM {
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
