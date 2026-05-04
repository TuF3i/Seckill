package dto

var (
	InvalidEmail        = Response{Status: 40001, Info: "Invalid Email"}
	InvalidPassword     = Response{Status: 40002, Info: "Invalid Password"}
	WrongPassword       = Response{Status: 40003, Info: "Wrong Password"}
	InvalidRefreshToken = Response{Status: 40004, Info: "Invalid Refresh Token"}
	WrongRefreshToken   = Response{Status: 40005, Info: "Wrong Refresh Token"}
	InvalidAccessToken  = Response{Status: 40006, Info: "Invalid Access Token"}
	WrongAccessToken    = Response{Status: 40007, Info: "Wrong Access Token"}
)

type Response struct {
	Status int32
	Info   string
}
