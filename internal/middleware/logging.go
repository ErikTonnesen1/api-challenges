package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		c.Next()

		log.Printf("%s: %s ~ [%d ms, || %d mu]",
			c.Request.Method,
			c.Request.URL,
			time.Since(start).Milliseconds(),
			time.Since(start).Microseconds(),
		)
	}
}
