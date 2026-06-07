package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/model"
	"secureops/backend-go/api/response"
)

type AssetController struct {
	assetService AssetService
}

func NewAssetController(assetService AssetService) *AssetController {
	return &AssetController{assetService: assetService}
}

func (c *AssetController) GetAssets(ctx *gin.Context) {
	assets, err := c.assetService.GetAllAssets(ctx.Request.Context())
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, assets)
}

func (c *AssetController) GetAsset(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.GetAsset(ctx.Request.Context(), id)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) CreateAsset(ctx *gin.Context) {
	var request model.AssetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.CreateAsset(ctx.Request.Context(), request)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, asset)
}

func (c *AssetController) UpdateAsset(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	var request model.AssetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.UpdateAsset(ctx.Request.Context(), id, request)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) DeleteAsset(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.DeleteAsset(ctx.Request.Context(), id)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) AssignVulnerability(ctx *gin.Context) {
	assetID, vulnerabilityID, ok := parsePair(ctx)
	if !ok {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.AssignVulnerability(ctx.Request.Context(), assetID, vulnerabilityID)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) RemoveVulnerability(ctx *gin.Context) {
	assetID, vulnerabilityID, ok := parsePair(ctx)
	if !ok {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.RemoveVulnerability(ctx.Request.Context(), assetID, vulnerabilityID)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func (c *AssetController) CalculateRisk(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		response.HandleGinError(ctx, response.ErrBadRequest)
		return
	}
	asset, err := c.assetService.CalculateRisk(ctx.Request.Context(), id)
	if err != nil {
		response.HandleGinError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, asset)
}

func parsePair(ctx *gin.Context) (int64, int64, bool) {
	assetID, err := parseID(ctx.Param("id"))
	if err != nil {
		return 0, 0, false
	}
	vulnerabilityID, err := parseID(ctx.Param("vulnerabilityId"))
	if err != nil {
		return 0, 0, false
	}
	return assetID, vulnerabilityID, true
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, strconv.ErrSyntax
	}
	return id, nil
}
