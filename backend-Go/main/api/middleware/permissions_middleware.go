package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/model"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("userRole")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": ErrForbidden.Message})
			return
		}

		value, ok := role.(string)
		if !ok || value != model.RoleAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": ErrForbidden.Message})
			return
		}

		ctx.Next()
	}
}
