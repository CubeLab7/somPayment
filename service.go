package somPayment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	config *Config
}

const (
	initiatePay   = "/v1/processing/init"
	postCheck     = "/v1/processing/"
	refund        = "/v1/processing/refund"
	exchangeRate  = "/v1/processing/exchangeRate"
	recurringList = "/v1/processing/recurringList"
)

func New(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) CartInit(ctx context.Context, data CartInitReq) (respBody []byte, response *InitPaymentResp, err error) {
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

	if respBody, err = sendRequest(s.config, inputs); err != nil {
		return
	}

	return
}

func (s *Service) Callback(ctx context.Context, data string) (respBody []byte, response CallbackResp, err error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	resp, err := decryptAES(encryptedBytes, []byte(s.config.Key))
	if err != nil {
		return
	}

	cleaned := cleanJSONString(resp)
	respBody = []byte(cleaned)

	if err = json.Unmarshal(respBody, &response); err != nil {
		return
	}

	return
}

// PostCheck - Получение статуса заказа
func (s *Service) PostCheck(ctx context.Context, orderID string) (respBody []byte, response *PostCheckResp, err error) {
	response = new(PostCheckResp)

	inputs := SendParams{
		Path:       postCheck,
		HttpMethod: http.MethodGet,
		QueryParams: map[string]string{
			"orderId": orderID,
		},
		Response: response,
	}

	if respBody, err = sendRequest(s.config, inputs); err != nil {
		return
	}

	return
}

// Refund - Проведениe полного возврата
func (s *Service) Refund(ctx context.Context, orderID string) (respBody []byte, response *RefundResp, err error) {
	response = new(RefundResp)

	inputs := SendParams{
		Path:       refund,
		HttpMethod: http.MethodDelete,
		QueryParams: map[string]string{
			"transactionId": orderID,
		},
		Response: response,
	}

	if respBody, err = sendRequest(s.config, inputs); err != nil {
		return
	}

	return
}

// ExchangeRate - Получение курса обмена валюты
func (s *Service) ExchangeRate(ctx context.Context) (respBody []byte, response *ExchangeRateResp, err error) {
	response = new(ExchangeRateResp)

	inputs := SendParams{
		Path:       exchangeRate,
		HttpMethod: http.MethodGet,
		Response:   response,
	}

	if respBody, err = sendRequest(s.config, inputs); err != nil {
		return
	}

	return
}

// RecurringList - Получения списка рекуррентов
func (s *Service) RecurringList(ctx context.Context) (respBody []byte, response *RecurringListResp, err error) {
	response = new(RecurringListResp)

	inputs := SendParams{
		Path:       recurringList,
		HttpMethod: http.MethodGet,
		Response:   response,
	}

	if respBody, err = sendRequest(s.config, inputs); err != nil {
		return
	}

	return
}

func sendRequest(config *Config, inputs SendParams) (respBody []byte, err error) {
	baseURL, err := url.Parse(config.URI)
	if err != nil {
		return respBody, fmt.Errorf("can't parse URI from config: %w", err)
	}

	// Добавляем путь из inputs.Path к базовому URL
	baseURL.Path += inputs.Path

	// Устанавливаем параметры запроса из queryParams
	query := baseURL.Query()
	for key, value := range inputs.QueryParams {
		query.Set(key, value)
	}
	baseURL.RawQuery = query.Encode()

	finalUrl := baseURL.String()

	req, err := http.NewRequest(inputs.HttpMethod, finalUrl, inputs.Body)
	if err != nil {
		return respBody, fmt.Errorf("can't create request for Som payment system! Err: %s", err)
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
		return respBody, fmt.Errorf("can't do request! Err: %s", err)
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return respBody, fmt.Errorf("can't read response body! Err: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("error: %v", string(respBody))
	}

	if err = json.Unmarshal(respBody, &inputs.Response); err != nil {
		return respBody, fmt.Errorf("can't unmarshall SomPayments resp: '%v'. Err: %w", string(respBody), err)
	}

	return
}
