package responses

import "net/http"

type BaseResponse struct {
	StatusText string `json:"statusText"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	IsSuccess  bool   `json:"isSuccess"`
}

func NewBaseResponse(statusCode int, message string, isSuccess bool) *BaseResponse {
	return &BaseResponse{
		StatusText: http.StatusText(statusCode),
		StatusCode: statusCode,
		Message:    message,
		IsSuccess:  isSuccess,
	}
}
