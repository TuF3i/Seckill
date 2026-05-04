package payment

type ProcessPaymentReq struct {
	OrderId string `json:"orderId"`
	UserId  string `json:"userId"`
}

type ProcessPaymentResp struct {
	Success bool `json:"success"`
}
