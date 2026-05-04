package payment

type ProcessPaymentReq struct {
	OrderId string `json:"orderId"`
}

type ProcessPaymentResp struct {
	Success bool `json:"success"`
}
