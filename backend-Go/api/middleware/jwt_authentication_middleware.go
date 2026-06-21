// Package middleware provides Gin middleware for request context setup, security guards, and request validation.
// JWTAuthenticationFilter validates incoming bearer tokens and establishes authenticated request state.
// It translates a verified JWT into typed identity data stored on GinContext.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/model"
	"secureops/backend-go/api/security"
)

// UserLookup defines how authentication middleware resolves a username to a user record.
// It accepts a request-scoped GinContext so lookup implementations can use the current request metadata.
type UserLookup interface {
	ExistsByUsername(ec *appcontext.GinContext, username string) (bool, error)
	FindByUsername(ec *appcontext.GinContext, username string) (model.User, error)
}

// JWTAuthenticationFilter validates Authorization bearer tokens, resolves the authenticated user,
// and stores typed authentication state on request context. It fails closed for missing, invalid,
// or unverifiable authentication.
func JWTAuthenticationFilter(jwtManager *security.JWTManager, users UserLookup) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			JWTAuthenticationEntryPoint(ctx)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		username, err := jwtManager.ExtractUsername(token)
		if err != nil {
			JWTAuthenticationEntryPoint(ctx)
			return
		}

		ec := appcontext.FromGinContext(ctx)
		exists, err := users.ExistsByUsername(ec, username)
		if err != nil || !exists {
			JWTAuthenticationEntryPoint(ctx)
			return
		}

		user, err := users.FindByUsername(ec, username)
		if err != nil {
			JWTAuthenticationEntryPoint(ctx)
			return
		}

		ec.SetUsername(username)
		ec.SetUserID(user.ID)
		ec.SetUserRole(user.Role)
		ctx.Next()
	}
}

// JWTAuthenticationEntryPoint aborts the request with a standard 401 Unauthorized response.
func JWTAuthenticationEntryPoint(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
