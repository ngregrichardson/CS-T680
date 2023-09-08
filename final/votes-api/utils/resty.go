package utils

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
)

func GetRequestError(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		var errorResponse ErrorResponse
		err := json.Unmarshal(response.Body(), &errorResponse)

		if err == nil {
			return errors.New(errorResponse.Error)
		} else {
			return errors.New("unknown error")
		}
	}

	return nil
}
