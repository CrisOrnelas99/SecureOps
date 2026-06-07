package controller

import (
	"encoding/json"
	"net/http"

	"risk-service-go/api/model"
	"risk-service-go/api/response"
	"risk-service-go/api/service"
)

type RiskController struct {
	riskService *service.RiskService
}

func NewRiskController(riskService *service.RiskService) *RiskController {
	return &RiskController{riskService: riskService}
}

func (c *RiskController) CalculateRisk(w http.ResponseWriter, r *http.Request) {
	var request model.RiskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, invalidRequestBody().Error())
		return
	}

	result, err := c.riskService.Calculate(request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, result)
}
