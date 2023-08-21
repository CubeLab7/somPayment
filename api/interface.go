package api

import (
	"context"
	"github.com/dwnGnL/somPayment/models"
)

type Api interface {
	CartInit(context.Context, models.RequestToSom) (models.ResponseFromSom, string, error)
	Callback(ctx context.Context, data string) (err error)
}
