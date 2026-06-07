package config

import (
	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/security"
)

func SecurityConfig(jwtService *security.JwtService, userLookup security.UserLookup) gin.HandlerFunc {
	return security.JwtAuthenticationFilter(jwtService, userLookup)
}
