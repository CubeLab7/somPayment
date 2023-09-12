package somPayment

import (
	"io"
)

type SendParams struct {
	Path       string
	HttpMethod string
	Body       io.Reader
	Response   interface{}
}

type CartInitReq struct {
	CurrencyCode int        `json:"currencyCode"`
	PayValue     int64      `json:"payValue"`
	Description  string     `json:"description"`
	SuccessLink  string     `json:"returnSuccessLink"`
	TimeToLive   int        `json:"timeToLive"`
	FailLink     string     `json:"returnFailLink"`
	CallbackUrl  string     `json:"callbackUrl"`
	Recurring    *Recurring `json:"recurring,omitempty"`
}

type Recurring struct {
	RecurringId string `json:"recurringId,omitempty"`
	ClientId    string `json:"clientId,omitempty"`
	ExpiryDate  string `json:"expiryDate,omitempty"`
	Frequency   int    `json:"frequency,omitempty"`
	Active      bool   `json:"active,omitempty"`
}

type InitPaymentResp struct {
	Code int          `json:"code"`
	Data CartInitResp `json:"data"`
}

type CartInitResp struct {
	Id           string  `json:"id"`
	ExchangeRate float64 `json:"exchangeRate"`
	PaySum       float64 `json:"paySum"`
	CurrencyCode int     `json:"currencyCode"`
	PayLink      string  `json:"payLink"`
	OrderId      string  `json:"orderId"`
}

type CallbackResp struct {
	OrderID    string `json:"orderId"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
	CreateDate string `json:"createDate"`
	UpdateDate string `json:"updateDate"`
	Recurring
	ProcessingStatus string `json:"processingStatus"`
	Pan              string `json:"pan"`
}
