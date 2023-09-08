package utils

import (
	"fmt"
	"polls-api/schema"

	"github.com/gin-gonic/gin"
)

func GenerateFullUrl(hostname string, path string) string {
	return hostname + path
}

func GetFullHostname(c *gin.Context) string {
	proto := "http"
	if c.Request.TLS != nil {
		proto = "https"
	}

	return proto + "://" + c.Request.Host
}

func FormatHostname(host string, port uint, tls bool) string {
	proto := "http"
	if tls {
		proto = "https"
	}

	return fmt.Sprintf("%s://%s:%d", proto, host, port)
}

func GenerateCRUDLinks(fullUrl string) gin.H {
	return gin.H{
		"get": schema.Link{
			Method: "GET",
			Url:    fullUrl,
		},
		"update": schema.Link{
			Method: "PATCH",
			Url:    fullUrl,
		},
		"delete": schema.Link{
			Method: "DELETE",
			Url:    fullUrl,
		},
	}
}
