package security

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserLookup interface {
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

func JwtAuthenticationFilter(jwtService *JwtService, users UserLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			JwtAuthenticationEntryPoint(c)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		username, err := jwtService.ExtractUsername(token)
		if err != nil {
			JwtAuthenticationEntryPoint(c)
			return
		}

		exists, err := users.ExistsByUsername(c.Request.Context(), username)
		if err != nil || !exists {
			JwtAuthenticationEntryPoint(c)
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
