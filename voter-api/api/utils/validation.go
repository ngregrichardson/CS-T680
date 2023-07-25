package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"voter-api/api/responses"

	"github.com/gin-gonic/gin"
)

func ValidateIntParam(c *gin.Context, param string) (uint, error) {
	stringValue := c.Param(param)

	v, err := strconv.ParseUint(stringValue, 10, 32)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewHttpBadRequestResponseWithMessage(fmt.Sprintf("Invalid value for %s. Expected integer.", param)))
		return 0, err
	}

	return uint(v), nil
}

func ValidateOptionalIntQuery(c *gin.Context, query string) (int, error) {
	stringValue := c.Query(query)

	if stringValue == "" {
		return 0, nil
	}

	v, err := strconv.ParseUint(stringValue, 10, 32)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewHttpBadRequestResponseWithMessage(fmt.Sprintf("Invalid value for %s. Expected integer.", query)))
		return 0, err
	}

	return int(v), nil
}
