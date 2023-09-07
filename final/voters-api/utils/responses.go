package utils

import "github.com/gin-gonic/gin"

func ResponseError(message string) gin.H {
	return gin.H{
		"error": message,
	}
}

func ResponseMessage(message string) gin.H {
	return gin.H{
		"message": message,
	}
}
