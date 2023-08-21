package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dwnGnL/somPayment/config"
	"github.com/dwnGnL/somPayment/lib"
	"github.com/dwnGnL/somPayment/models"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type Service struct {
	som    *lib.SOM
	config *config.Config
}

const (
	initiatePay = "/v1/processing/init"
)

func New(config *config.Config, som *lib.SOM) *Service {
	return &Service{
		som:    som,
		config: config,
	}
}

func (s *Service) CartInit(ctx context.Context, data models.RequestToSom) (response models.ResponseFromSom, dataSave string, err error) {
	// отправка в SOM
	body := new(bytes.Buffer)
	err = jsoniter.NewEncoder(body).Encode(data)
	if err != nil {
		err = fmt.Errorf("can't encode request: %s", err)
		return
	}

	dataSave = fmt.Sprint(body)

	req, err := s.som.PrepareRequest(ctx, initiatePay, body)
	if err != nil {
		return
	}

	response, err = s.som.SendRequest(req)
	if err != nil {
		return
	}

	return
}

func (s *Service) Callback(ctx context.Context, data string) (err error) {
	// Декодируем тело из Base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	resp, err := decryptAES(encryptedBytes, []byte(s.config.Key))
	if err != nil {
		return
	}

	cleaned := strings.ReplaceAll(string(resp), "\u0001", "")

	var response CallbackReq
	if err = jsoniter.Unmarshal([]byte(cleaned), &response); err != nil {
		return
	}

	return
}
