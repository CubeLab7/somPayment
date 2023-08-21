package models

import "github.com/google/uuid"

type RequestToSom struct {
	CurrencyCode      int       `json:"currencyCode"`
	PayValue          int64     `json:"payValue"`
	Description       string    `json:"description"`
	ReturnSuccessLink string    `json:"returnSuccessLink"`
	ReturnFailLink    string    `json:"returnFailLink"`
	CallbackUrl       string    `json:"callbackUrl"`
	Recurring         Recurring `json:"recurring,omitempty"`
}

type Recurring struct {
	ClientId   uuid.UUID `json:"clientId"`
	ExpiryDate string    `json:"expiryDate"`
	Frequency  int       `json:"frequency"`
	Active     bool      `json:"active"`
}

type ResponseFromSom struct {
	Code int          `json:"code"`
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Id           string  `json:"id"`
	ExchangeRate float64 `json:"exchangeRate"`
	PaySum       float64 `json:"paySum"`
	CurrencyCode int     `json:"currencyCode"`
	PayLink      string  `json:"payLink"`
	OrderId      string  `json:"orderId"`
}
