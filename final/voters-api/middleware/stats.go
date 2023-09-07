package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

var (
	BootTime         time.Time
	RequestCount     int
	TotalRequestTime time.Duration
	ErrorCount       int
)

func Stats() gin.HandlerFunc {
	BootTime = time.Now()
	return func(c *gin.Context) {
		RequestCount++
		start := time.Now()
		c.Next()
		TotalRequestTime += time.Since(start)

		if c.Writer.Status() >= 400 {
			ErrorCount++
		}
	}
}

func GetStats() gin.H {
	return gin.H{
		"uptime":             time.Since(BootTime).String(),
		"requests":           RequestCount,
		"totalRequestTime":   TotalRequestTime.String(),
		"averageRequestTime": time.Duration(float64(TotalRequestTime) / float64(RequestCount)).String(),
		"errors":             ErrorCount,
	}
}
