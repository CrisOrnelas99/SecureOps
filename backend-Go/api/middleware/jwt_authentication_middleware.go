package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/model"
	"secureops/backend-go/api/security"
)

type UserLookup interface {
	ExistsByUsername(ec *appcontext.GinContext, username string) (bool, error)
	FindByUsername(ec *appcontext.GinContext, username string) (model.User, error)
}

func JwtAuthenticationFilter(jwtManager *security.JWTManager, users UserLookup) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			JwtAuthenticationEntryPoint(ctx)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		username, err := jwtManager.ExtractUsername(token)
		if err != nil {
			JwtAuthenticationEntryPoint(ctx)
			return
		}

		ec := appcontext.FromGinContext(ctx)
		exists, err := users.ExistsByUsername(ec, username)
		if err != nil || !exists {
			JwtAuthenticationEntryPoint(ctx)
			return
		}

		user, err := users.FindByUsername(ec, username)
		if err != nil {
			JwtAuthenticationEntryPoint(ctx)
			return
		}

		ctx.Set("username", username)
		ctx.Set("userID", user.ID)
		ctx.Set("userRole", user.Role)
		ctx.Next()
	}
}

func JwtAuthenticationEntryPoint(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
