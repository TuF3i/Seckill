package lkeygen

func GenOrderStatusKey(orderId string) string {
	return "order:" + "status:" + orderId
}

func GenOrderInfoKey(orderId string) string {
	return "order:" + "info:" + orderId
}

func GenUserOrdersKey(userId string) string {
	return "order:" + "user:" + userId
}
