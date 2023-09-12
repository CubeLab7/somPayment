package somPayment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	inputs := SendParams{
		Path:       initiatePay,
		HttpMethod: http.MethodPost,
		Body:       body,
		Response:   response,
	}

	if err = sendRequest(s.config, inputs); err != nil {
		return
	}

	return
}

func (s *Service) Callback(ctx context.Context, data string) (response CallbackResp, err error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	resp, err := decryptAES(encryptedBytes, []byte(s.config.Key))
	if err != nil {
		return
	}

	cleaned := cleanJSONString(resp)

	if err = json.Unmarshal([]byte(cleaned), &response); err != nil {
		return
	}

	return
}

func sendRequest(config *Config, inputs SendParams) (err error) {
	finalUrl := fmt.Sprintf("%v/%v", config.URI, inputs.Path)

	req, err := http.NewRequest(inputs.HttpMethod, finalUrl, inputs.Body)
	if err != nil {
		return fmt.Errorf("can't create request for Som payment system! Err: %s", err)
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
		return fmt.Errorf("can't do request! Err: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read response body! Err: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %v", string(respBody))
	}

	if err = json.Unmarshal(respBody, &inputs.Response); err != nil {
		return fmt.Errorf("can't unmarshall SomPayments resp: '%v'. Err: %w", string(respBody), err)
	}

	return
}
