// Package middleware provides Gin middleware for request context setup, security guards, and request validation.
// CORS middleware enforces explicit origin allowlisting and handles preflight requests safely.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors returns middleware that adds CORS headers and handles OPTIONS preflight requests.
func Cors(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
