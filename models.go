package somPayment

import (
	"github.com/google/uuid"
	"io"
)

type SendParams struct {
	Path       string
	HttpMethod string
	Body       io.Reader
	Response   interface{}
}

type CartInitReq struct {
	CurrencyCode      int        `json:"currencyCode"`
	PayValue          int64      `json:"payValue"`
	Description       string     `json:"description"`
	ReturnSuccessLink string     `json:"returnSuccessLink"`
	TimeToLive        int        `json:"timeToLive"`
	ReturnFailLink    string     `json:"returnFailLink"`
	CallbackUrl       string     `json:"callbackUrl"`
	Recurring         *Recurring `json:"recurring,omitempty"`
}

type Recurring struct {
	RecurringId string    `json:"recurringId"`
	ClientId    uuid.UUID `json:"clientId"`
	ExpiryDate  string    `json:"expiryDate"`
	Frequency   int       `json:"frequency"`
	Active      bool      `json:"active"`
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
	OrderID          string    `json:"orderId"`
	Status           int       `json:"status"`
	StatusName       string    `json:"statusName"`
	CreateDate       string    `json:"createDate"`
	UpdateDate       string    `json:"updateDate"`
	Recurring        Recurring `json:"recurring"`
	ProcessingStatus string    `json:"processingStatus"`
	Pan              string    `json:"pan"`
}
