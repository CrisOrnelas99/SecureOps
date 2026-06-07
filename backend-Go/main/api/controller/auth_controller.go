package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/model"
	"secureops/backend-go/api/response"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var request model.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}

	if err := c.authService.Register(ctx.Request.Context(), request); err != nil {
		response.HandleGinError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request model.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}

	token, err := c.authService.Login(ctx.Request.Context(), request)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, model.LoginResponse{Token: token})
}
