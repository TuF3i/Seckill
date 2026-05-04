package dto

var (
	InvalidOrderID    = Response{Status: 43001, Info: "Invalid Order ID"}
	OrderNotFound     = Response{Status: 43002, Info: "Order Not Found"}
	OrderAlreadyPaid  = Response{Status: 43003, Info: "Order Already Paid"}
	PaymentFailed     = Response{Status: 43004, Info: "Payment Processing Failed"}
)

type Response struct {
	Status int32
	Info   string
}
