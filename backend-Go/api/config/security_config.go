package config

import (
	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/middleware"
	"secureops/backend-go/api/security"
)

func SecurityConfig(jwtManager *security.JWTManager, userLookup middleware.UserLookup) gin.HandlerFunc {
	return middleware.JWTAuthenticationFilter(jwtManager, userLookup)
}
