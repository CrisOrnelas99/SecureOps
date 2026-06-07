package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"secureops/backend-go/api/model"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) FindAll(ctx context.Context) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.WithContext(ctx).Order("id").Find(&assets).Error
	return assets, err
}

func (r *AssetRepository) FindByID(ctx context.Context, id int64) (model.Asset, error) {
	var asset model.Asset
	err := r.db.WithContext(ctx).Preload("Vulnerabilities").First(&asset, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Asset{}, ErrAssetNotFound
	}
	return asset, err
}

func (r *AssetRepository) Save(ctx context.Context, request model.AssetRequest) (model.Asset, error) {
	asset := model.Asset{
		Name:            request.Name,
		Type:            request.Type,
		IPAddress:       request.IPAddress,
		OperatingSystem: request.OperatingSystem,
		Owner:           request.Owner,
		Criticality:     request.Criticality,
		RiskScore:       0,
		RiskLevel:       "Low",
	}
	err := r.db.WithContext(ctx).Create(&asset).Error
	return asset, err
}

func (r *AssetRepository) Update(ctx context.Context, id int64, request model.AssetRequest) (model.Asset, error) {
	asset, err := r.FindByID(ctx, id)
	if err != nil {
		return model.Asset{}, err
	}

	asset.Name = request.Name
	asset.Type = request.Type
	asset.IPAddress = request.IPAddress
	asset.OperatingSystem = request.OperatingSystem
	asset.Owner = request.Owner
	asset.Criticality = request.Criticality

	err = r.db.WithContext(ctx).Save(&asset).Error
	if err != nil {
		return model.Asset{}, err
	}
	return r.FindByID(ctx, id)
}

func (r *AssetRepository) Delete(ctx context.Context, id int64) (model.Asset, error) {
	asset, err := r.FindByID(ctx, id)
	if err != nil {
		return model.Asset{}, err
	}
	err = r.db.WithContext(ctx).Delete(&asset).Error
	return asset, err
}

func (r *AssetRepository) AssignVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error) {
	asset, err := r.FindByID(ctx, assetID)
	if err != nil {
		return model.Asset{}, err
	}

	var vulnerability model.Vulnerability
	err = r.db.WithContext(ctx).First(&vulnerability, vulnerabilityID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Asset{}, ErrVulnerabilityNotFound
	}
	if err != nil {
		return model.Asset{}, err
	}

	for _, assigned := range asset.Vulnerabilities {
		if assigned.ID == vulnerability.ID {
			return model.Asset{}, ErrDuplicateAssignment
		}
	}

	err = r.db.WithContext(ctx).Model(&asset).Association("Vulnerabilities").Append(&vulnerability)
	if err != nil {
		return model.Asset{}, err
	}

	return r.FindByID(ctx, assetID)
}

func (r *AssetRepository) RemoveVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error) {
	asset, vulnerability, err := r.findAssetAndVulnerability(ctx, assetID, vulnerabilityID)
	if err != nil {
		return model.Asset{}, err
	}

	err = r.db.WithContext(ctx).Model(&asset).Association("Vulnerabilities").Delete(&vulnerability)
	if err != nil {
		return model.Asset{}, err
	}

	return r.FindByID(ctx, assetID)
}

func (r *AssetRepository) FindVulnerabilities(ctx context.Context, assetID int64) ([]model.Vulnerability, error) {
	asset, err := r.FindByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return asset.Vulnerabilities, nil
}

func (r *AssetRepository) PersistRiskResult(ctx context.Context, id int64, riskResponse model.RiskCalculationResponse) (model.Asset, error) {
	if riskResponse.RiskScore < -32768 || riskResponse.RiskScore > 32767 {
		return model.Asset{}, ErrRiskScoreOutOfRange
	}

	asset, err := r.FindByID(ctx, id)
	if err != nil {
		return model.Asset{}, err
	}

	asset.RiskScore = int16(riskResponse.RiskScore)
	asset.RiskLevel = riskResponse.RiskLevel

	err = r.db.WithContext(ctx).Save(&asset).Error
	if err != nil {
		return model.Asset{}, err
	}

	return r.FindByID(ctx, id)
}

func (r *AssetRepository) EnsureAssetAndVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) error {
	_, _, err := r.findAssetAndVulnerability(ctx, assetID, vulnerabilityID)
	return err
}

func (r *AssetRepository) findAssetAndVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, model.Vulnerability, error) {
	asset, err := r.FindByID(ctx, assetID)
	if err != nil {
		return model.Asset{}, model.Vulnerability{}, err
	}

	var vulnerability model.Vulnerability
	err = r.db.WithContext(ctx).First(&vulnerability, vulnerabilityID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Asset{}, model.Vulnerability{}, ErrVulnerabilityNotFound
	}
	if err != nil {
		return model.Asset{}, model.Vulnerability{}, err
	}

	return asset, vulnerability, nil
}
