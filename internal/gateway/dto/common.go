package dto

var (
	EmptyJWTString = Response{Status: 10001, Info: "Empty JWT String"} // EmptyJWTString 空JWT字符串
)

type FinalResponse struct {
	Status int32       `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Status int32  `json:"status"`
	Info   string `json:"info"`
}

func (r Response) Error() string {
	return r.Info
}

func InternalError(err error) Response {
	return Response{
		Status: 500,
		Info:   err.Error(),
	}
}

func GenFinalResponse(response Response, data interface{}) FinalResponse {
	return FinalResponse{Status: response.Status, Info: response.Info, Data: data}
}

func GenBizErrorResponse(status int32, info string) Response {
	return Response{Status: status, Info: info}
}
