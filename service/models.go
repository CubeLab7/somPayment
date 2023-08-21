package service

type CallbackReq struct {
	OrderID          string `json:"orderId"`
	Status           int    `json:"status"`
	StatusName       string `json:"statusName"`
	CreateDate       string `json:"createDate"`
	UpdateDate       string `json:"updateDate"`
	RecurringId      string `json:"recurringId"`
	ProcessingStatus string `json:"processingStatus"`
	Pan              string `json:"pan"`
}
