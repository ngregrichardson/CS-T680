package utils

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func ResponseError(message string) ErrorResponse {
	return ErrorResponse{
		Error: message,
	}
}

func ResponseMessage(message string) MessageResponse {
	return MessageResponse{
		Message: message,
	}
}
