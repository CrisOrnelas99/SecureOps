package service

import (
	"context"
	"net"
	"strings"

	"secureops/backend-go/api/model"
)

type AssetService struct {
	assetRepository         AssetRepository
	vulnerabilityRepository VulnerabilityRepository
	restClient              *RestClient
	assetRiskService        *AssetRiskService
}

func NewAssetService(
	assetRepository AssetRepository,
	vulnerabilityRepository VulnerabilityRepository,
	restClient *RestClient,
	assetRiskService *AssetRiskService,
) *AssetService {
	return &AssetService{
		assetRepository:         assetRepository,
		vulnerabilityRepository: vulnerabilityRepository,
		restClient:              restClient,
		assetRiskService:        assetRiskService,
	}
}

func (s *AssetService) GetAllAssets(ctx context.Context) ([]model.Asset, error) {
	assets, err := s.assetRepository.FindAll(ctx)
	return assets, mapRepositoryError(err)
}

func (s *AssetService) GetAsset(ctx context.Context, id int64) (model.Asset, error) {
	asset, err := s.assetRepository.FindByID(ctx, id)
	return asset, mapRepositoryError(err)
}

func (s *AssetService) CreateAsset(ctx context.Context, request model.AssetRequest) (model.Asset, error) {
	if err := validateAssetRequest(request); err != nil {
		return model.Asset{}, err
	}
	asset, err := s.assetRepository.Save(ctx, request)
	return asset, mapRepositoryError(err)
}

func (s *AssetService) UpdateAsset(ctx context.Context, id int64, request model.AssetRequest) (model.Asset, error) {
	if err := validateAssetRequest(request); err != nil {
		return model.Asset{}, err
	}
	asset, err := s.assetRepository.Update(ctx, id, request)
	return asset, mapRepositoryError(err)
}

func (s *AssetService) DeleteAsset(ctx context.Context, id int64) (model.Asset, error) {
	asset, err := s.assetRepository.Delete(ctx, id)
	return asset, mapRepositoryError(err)
}

func (s *AssetService) AssignVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error) {
	asset, err := s.assetRepository.AssignVulnerability(ctx, assetID, vulnerabilityID)
	return asset, mapRepositoryError(err)
}

func (s *AssetService) RemoveVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error) {
	asset, err := s.assetRepository.RemoveVulnerability(ctx, assetID, vulnerabilityID)
	return asset, mapRepositoryError(err)
}

func (s *AssetService) CalculateRisk(ctx context.Context, id int64) (model.Asset, error) {
	request, err := s.assetRiskService.LoadRiskCalculationRequest(ctx, id)
	if err != nil {
		return model.Asset{}, mapRepositoryError(err)
	}

	response, err := s.restClient.CalculateRisk(ctx, request)
	if err != nil {
		return model.Asset{}, err
	}

	asset, err := s.assetRiskService.PersistRiskResult(ctx, id, response)
	return asset, mapRepositoryError(err)
}

func validateAssetRequest(request model.AssetRequest) error {
	if strings.TrimSpace(request.Name) == "" ||
		strings.TrimSpace(request.Type) == "" ||
		strings.TrimSpace(request.Owner) == "" ||
		strings.TrimSpace(request.Criticality) == "" {
		return ErrInvalidRequestData
	}

	if ip := net.ParseIP(request.IPAddress); ip == nil || ip.To4() == nil {
		return ErrInvalidRequestData
	}

	return nil
}
