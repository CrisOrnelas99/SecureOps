// Package middleware provides Gin middleware for request context setup, security guards, and request validation.
// GormMiddleware attaches the shared GORM database handle into request-scoped context.
package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	appcontext "secureops/backend-go/api/context"
)

// GormMiddleware attaches the shared GORM database handle to the request context.
// This middleware does not open or close the database connection; it simply makes the handle available
// to downstream handlers and services via appcontext.GinContext.
func GormMiddleware(database *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ec := appcontext.FromGinContext(ctx)
		ec.SetDatabase(database)
		ctx.Next()
	}
}
