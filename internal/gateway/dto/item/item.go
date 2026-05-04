package item

type AddItemReq struct {
	Name        string  `json:"name"`
	Stock       int64   `json:"stock"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type AddItemResp struct {
	ItemId string `json:"itemId"`
}

type DeleteItemReq struct {
	ItemId string `json:"itemId"`
}

type FlashSaleReq struct {
	ItemId string `json:"itemId"`
}

type ListItemsResp struct {
	Items []ItemInfoResp `json:"items"`
}

type ItemInfoResp struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Stock       int64   `json:"stock"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
