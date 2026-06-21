// Package service provides asset-related application services.
package service

import (
	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/model"
	baserepository "secureops/backend-go/api/repository"
	baseservice "secureops/backend-go/api/service"
)

type assetServiceImpl struct {
	assetRepository baserepository.AssetRepository
}

// NewAssetService creates an asset service backed by the supplied repository.
func NewAssetService(assetRepository baserepository.AssetRepository) baseservice.AssetService {
	return &assetServiceImpl{assetRepository: assetRepository}
}

// GetAllAssets returns all assets owned by the authenticated user.
func (s *assetServiceImpl) GetAllAssets(ec *appcontext.GinContext) ([]model.Asset, error) {
	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return nil, err
	}
	assets, err := s.assetRepository.FindAllByUser(ec, userID)
	return assets, baseservice.TranslateRepositoryError(err)
}

// GetAsset returns a single asset owned by the authenticated user.
func (s *assetServiceImpl) GetAsset(ec *appcontext.GinContext, id int64) (model.Asset, error) {
	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return model.Asset{}, err
	}
	asset, err := s.assetRepository.FindByIDForUser(ec, id, userID)
	return asset, baseservice.TranslateRepositoryError(err)
}

// CreateAsset validates and saves a new asset for the authenticated user.
func (s *assetServiceImpl) CreateAsset(ec *appcontext.GinContext, asset model.Asset) (model.Asset, error) {
	if err := baseservice.ValidateAsset(asset); err != nil {
		return model.Asset{}, err
	}

	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return model.Asset{}, err
	}
	asset.UserID = userID

	created, err := s.assetRepository.Save(ec, asset)
	return created, baseservice.TranslateRepositoryError(err)
}

// UpdateAsset validates and updates an existing asset for the authenticated user.
func (s *assetServiceImpl) UpdateAsset(ec *appcontext.GinContext, id int64, asset model.Asset) (model.Asset, error) {
	if err := baseservice.ValidateAsset(asset); err != nil {
		return model.Asset{}, err
	}

	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return model.Asset{}, err
	}

	updated, err := s.assetRepository.UpdateForUser(ec, id, userID, asset)
	return updated, baseservice.TranslateRepositoryError(err)
}

// DeleteAsset removes an asset owned by the authenticated user.
func (s *assetServiceImpl) DeleteAsset(ec *appcontext.GinContext, id int64) (model.Asset, error) {
	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return model.Asset{}, err
	}
	asset, err := s.assetRepository.DeleteForUser(ec, id, userID)
	return asset, baseservice.TranslateRepositoryError(err)
}

// AssignVulnerability attaches a vulnerability to an asset owned by the authenticated user.
func (s *assetServiceImpl) AssignVulnerability(ec *appcontext.GinContext, assetID int64, vulnerabilityID int64) (model.Asset, error) {
	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return model.Asset{}, err
	}
	asset, err := s.assetRepository.AssignVulnerabilityForUser(ec, assetID, userID, vulnerabilityID)
	return asset, baseservice.TranslateRepositoryError(err)
}

// RemoveVulnerability removes a vulnerability from an asset owned by the authenticated user.
func (s *assetServiceImpl) RemoveVulnerability(ec *appcontext.GinContext, assetID int64, vulnerabilityID int64) (model.Asset, error) {
	userID, err := baseservice.AuthenticatedUserID(ec)
	if err != nil {
		return model.Asset{}, err
	}
	asset, err := s.assetRepository.RemoveVulnerabilityForUser(ec, assetID, userID, vulnerabilityID)
	return asset, baseservice.TranslateRepositoryError(err)
}
