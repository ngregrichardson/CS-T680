package responses

import "net/http"

type DataResponse struct {
	*BaseResponse
	Data interface{} `json:"data"`
}

func NewDataResponse(statusCode int, message string, data interface{}) *DataResponse {
	return &DataResponse{
		BaseResponse: NewBaseResponse(statusCode, message, true),
		Data:         data,
	}
}

func NewHttpOkResponseWithMessage(message string, data interface{}) *DataResponse {
	return NewDataResponse(http.StatusOK, message, data)
}

func NewHttpOkResponse(data interface{}) *DataResponse {
	return NewHttpOkResponseWithMessage("Success", data)
}

func NewHttpCreatedResponseWithMessage(message string, data interface{}) *DataResponse {
	return NewDataResponse(http.StatusCreated, message, data)
}

func NewHttpCreatedResponse(data interface{}) *DataResponse {
	return NewHttpCreatedResponseWithMessage("Success", data)
}
