package dto

type FinalResponse struct {
	Status int64       `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Status int64  `json:"status"`
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
