package service

import (
	"context"
	"strings"

	"secureops/backend-go/api/model"
)

type AssetRiskService struct {
	assetRepository AssetRepository
}

func NewAssetRiskService(assetRepository AssetRepository) *AssetRiskService {
	return &AssetRiskService{assetRepository: assetRepository}
}

func (s *AssetRiskService) LoadRiskCalculationRequest(ctx context.Context, id int64) (model.RiskCalculationRequest, error) {
	asset, err := s.assetRepository.FindByID(ctx, id)
	if err != nil {
		return model.RiskCalculationRequest{}, mapRepositoryError(err)
	}

	return buildRiskCalculationRequest(asset), nil
}

func (s *AssetRiskService) PersistRiskResult(ctx context.Context, id int64, response model.RiskCalculationResponse) (model.Asset, error) {
	asset, err := s.assetRepository.PersistRiskResult(ctx, id, response)
	return asset, mapRepositoryError(err)
}

func buildRiskCalculationRequest(asset model.Asset) model.RiskCalculationRequest {
	request := model.RiskCalculationRequest{
		AssetID:     asset.ID,
		Criticality: asset.Criticality,
	}

	for _, vulnerability := range asset.Vulnerabilities {
		switch strings.ToLower(strings.TrimSpace(vulnerability.Severity)) {
		case "critical":
			request.CriticalVulnerabilities++
		case "high":
			request.HighVulnerabilities++
		case "medium":
			request.MediumVulnerabilities++
		case "low":
			request.LowVulnerabilities++
		}
	}

	return request
}
