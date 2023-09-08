package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateIntParam(c *gin.Context, param string) (uint, error) {
	stringValue := c.Param(param)

	v, err := strconv.ParseUint(stringValue, 10, 32)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseError(fmt.Sprintf("invalid value for %s. Expected integer.", param)))
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
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseError(fmt.Sprintf("invalid value for %s. Expected integer.", query)))
		return 0, err
	}

	return int(v), nil
}

func ValidateBody(c *gin.Context, body interface{}) error {
	if err := c.ShouldBindJSON(body); err != nil {
		actualBody, _ := json.Marshal(body)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseError(fmt.Sprintf("invalid body format. follow the following format: %s", string(actualBody))))
		return err
	}

	return nil
}
