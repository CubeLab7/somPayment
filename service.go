package somPayment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

func (s *Service) CartInit(ctx context.Context, data CartInitReq) (response *InitPaymentResp, err error) {
	response = new(InitPaymentResp)

	// отправка в SOM
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(data)
	if err != nil {
		err = fmt.Errorf("can't encode request: %s", err)
		return
	}

	if err = sendRequest(initiatePay, body, s.config, response); err != nil {
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

func sendRequest(method string, body io.Reader, config *Config, response interface{}) (err error) {
	url := fmt.Sprintf("%v/%v", config.URI, method)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		err = fmt.Errorf("can't create request for Som payment system: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", basicAuth(config.Login, config.Pass))

	httpClient := http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Second * time.Duration(config.IdleConnTimeoutSec),
		},
		Timeout: time.Second * time.Duration(config.RequestTimeoutSec),
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("can't do request: %s", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("can't read response body: %s", err)
		return
	}

	log.Println("Resp: ", string(respBody))

	if err = json.Unmarshal(respBody, &response); err != nil {
		err = fmt.Errorf("can't unmarshall SomPayments resp: '%v'. Err: %w", string(respBody), err)
		return
	}

	return
}
