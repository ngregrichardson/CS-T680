package responses

import "net/http"

type ErrorResponse struct {
	*BaseResponse
}

func NewErrorResponse(statusCode int, message string) *ErrorResponse {
	return &ErrorResponse{
		BaseResponse: NewBaseResponse(statusCode, message, false),
	}
}

func NewHttpBadRequestResponseWithMessage(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, message)
}

func NewHttpBadRequestResponse() *ErrorResponse {
	return NewHttpBadRequestResponseWithMessage("Bad Request")
}

func NewHttpNotFoundResponseWithMessage(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusNotFound, message)
}

func NewHttpNotFoundResponse(message string) *ErrorResponse {
	return NewHttpNotFoundResponseWithMessage("Not Found")
}

func NewHttpConflictResponseWithMessage(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusConflict, message)
}

func NewHttpConflictResponse(message string) *ErrorResponse {
	return NewHttpConflictResponseWithMessage("Conflict")
}

func NewHttpInternalServerErrorResponseWithMessage(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError, message)
}

func NewHttpInternalServerErrorResponse(message string) *ErrorResponse {
	return NewHttpInternalServerErrorResponseWithMessage("A unexpected issue occurred")
}
