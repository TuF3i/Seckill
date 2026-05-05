package models

type OrderMessage struct {
	OrderId string  `json:"orderId"`
	UserId  string  `json:"userId"`
	ItemId  string  `json:"itemId"`
	Price   float64 `json:"price"`
}
