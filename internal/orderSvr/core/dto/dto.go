package dto

var (
	InvalidOrderID     = Response{Status: 42001, Info: "Invalid Order ID"}
	InvalidUserID      = Response{Status: 42002, Info: "Invalid User ID"}
	InvalidItemID      = Response{Status: 42003, Info: "Invalid Item ID"}
	InvalidPrice       = Response{Status: 42004, Info: "Invalid Price"}
	OrderNotFound      = Response{Status: 42005, Info: "Order Not Found"}
	OrderAlreadyExists = Response{Status: 42006, Info: "Order Already Exists"}
	MQSendFailed       = Response{Status: 42007, Info: "MQ Send Failed"}
)

type Response struct {
	Status int32
	Info   string
}
