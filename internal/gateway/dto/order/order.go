package order

type CreateOrderReq struct {
	ItemId string `json:"itemId"`
}

type CreateOrderResp struct {
	OrderId string `json:"orderId"`
}

type QueryOrdersReq struct {
	UserId string `json:"userId"`
}

type OrderInfoResp struct {
	OrderId    string  `json:"orderId"`
	UserId     string  `json:"userId"`
	ItemId     string  `json:"itemId"`
	Price      float64 `json:"price"`
	Status     int32   `json:"status"`
	CreateTime string  `json:"createTime"`
}
