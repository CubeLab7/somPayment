package somPayment

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dwnGnL/somPayment/lib"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type Service struct {
	config *Config
}

const (
	initiatePay = "/v1/processing/init"
)

func New(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) CartInit(ctx context.Context, data CartInitReq) (response InitPaymentResp, err error) {
	// отправка в SOM
	body := new(bytes.Buffer)
	err = jsoniter.NewEncoder(body).Encode(data)
	if err != nil {
		err = fmt.Errorf("can't encode request: %s", err)
		return
	}

	request, err := lib.PrepareRequest(ctx, initiatePay, body)
	if err != nil {
		return
	}

	if err = lib.SendRequest(request, response); err != nil {
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
