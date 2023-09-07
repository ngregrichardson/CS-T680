package utils

import (
	"fmt"

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
