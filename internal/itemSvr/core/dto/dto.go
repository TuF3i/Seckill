package dto

var (
	InvalidItemName   = Response{Status: 41001, Info: "Invalid Item Name"}
	InvalidItemStock  = Response{Status: 41002, Info: "Invalid Item Stock"}
	InvalidItemPrice  = Response{Status: 41003, Info: "Invalid Item Price"}
	InvalidItemID     = Response{Status: 41004, Info: "Invalid Item ID"}
	ItemNotFound      = Response{Status: 41005, Info: "Item Not Found"}
	FlashAlreadyStart = Response{Status: 41006, Info: "Flash Sale Already Started"}
	FlashNotStart     = Response{Status: 41007, Info: "Flash Sale Not Started"}
	PermissionDenied  = Response{Status: 41008, Info: "Permission Denied"}
)

type Response struct {
	Status int32
	Info   string
}
