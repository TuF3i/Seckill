package lkeygen

func GenItemFlashKey(itemId string) string {
	return "item:" + "flash:" + itemId
}

func GenItemStockKey(itemId string) string {
	return "item:" + "stock:" + itemId
}

func GenItemPurchaseLimitKey(itemId string, userId string) string {
	return "item:" + "limit:" + itemId + ":" + userId
}
