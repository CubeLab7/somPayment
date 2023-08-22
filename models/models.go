package models

import "github.com/google/uuid"

type CartInit struct {
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

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
