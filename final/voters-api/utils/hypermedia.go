package utils

import (
	"fmt"
	"voters-api/schema"
)

func FormatHostname(port uint, tls bool) string {
	proto := "http"
	if tls {
		proto = "https"
	}

	host := GetEnvironmentVariable("VISIBLE_API_HOST", "localhost")

	return fmt.Sprintf("%s://%s:%d", proto, host, port)
}

func GenerateCRUDLinks(fullUrl string) schema.Links {
	return schema.Links{
		Get: schema.Link{
			Method: "GET",
			Url:    fullUrl,
		},
		Update: schema.Link{
			Method: "PATCH",
			Url:    fullUrl,
		},
		Delete: schema.Link{
			Method: "DELETE",
			Url:    fullUrl,
		},
	}
}
