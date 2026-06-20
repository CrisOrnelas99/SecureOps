package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/model"
	baserepository "secureops/backend-go/api/repository"
	"secureops/backend-go/api/utils"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) database(ec *appcontext.GinContext) *gorm.DB {
	if ec != nil && ec.Database() != nil {
		return ec.Database()
	}
	return r.db
}

func (r *AssetRepository) FindAllByUser(ec *appcontext.GinContext, userID int64) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.database(ec).WithContext(ec.RequestContext()).Where("user_id = ?", userID).Order("id").Find(&assets).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baserepository.ErrReadFailed, err)
	}
	return assets, nil
}

func (r *AssetRepository) FindByIDForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Asset, error) {
	var asset model.Asset
	err := r.database(ec).WithContext(ec.RequestContext()).
		Preload("Vulnerabilities", "user_id = ?", userID).
		Where("user_id = ?", userID).
		First(&asset, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Asset{}, baserepository.ErrAssetNotFound
	}
	if err != nil {
		return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrReadFailed, err)
	}
	return asset, nil
}

func (r *AssetRepository) Save(ec *appcontext.GinContext, asset model.Asset) (model.Asset, error) {
	if asset.UserID <= 0 || asset.Name == "" || asset.Type == "" || asset.IPAddress == "" || asset.Owner == "" || asset.Criticality == "" {
		return model.Asset{}, baserepository.ErrInvalidData
	}

	err := r.database(ec).WithContext(ec.RequestContext()).Create(&asset).Error
	if err != nil {
		databaseErr := utils.TranslateDatabaseError(err)
		if errors.Is(databaseErr, utils.ErrForeignKeyViolation) {
			return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrInvalidReference, databaseErr)
		}
		if errors.Is(databaseErr, utils.ErrCheckConstraintViolation) {
			return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrInvalidData, databaseErr)
		}
		return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrCreateFailed, databaseErr)
	}
	return asset, nil
}

func (r *AssetRepository) UpdateForUser(ec *appcontext.GinContext, id int64, userID int64, updates model.Asset) (model.Asset, error) {
	if updates.Name == "" || updates.Type == "" || updates.IPAddress == "" || updates.Owner == "" || updates.Criticality == "" {
		return model.Asset{}, baserepository.ErrInvalidData
	}

	asset, err := r.FindByIDForUser(ec, id, userID)
	if err != nil {
		return model.Asset{}, err
	}

	asset.Name = updates.Name
	asset.Type = updates.Type
	asset.IPAddress = updates.IPAddress
	asset.OperatingSystem = updates.OperatingSystem
	asset.Owner = updates.Owner
	asset.Criticality = updates.Criticality

	err = r.database(ec).WithContext(ec.RequestContext()).Save(&asset).Error
	if err != nil {
		databaseErr := utils.TranslateDatabaseError(err)
		if errors.Is(databaseErr, utils.ErrForeignKeyViolation) {
			return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrInvalidReference, databaseErr)
		}
		if errors.Is(databaseErr, utils.ErrCheckConstraintViolation) {
			return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrInvalidData, databaseErr)
		}
		return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrUpdateFailed, databaseErr)
	}
	return r.FindByIDForUser(ec, id, userID)
}

func (r *AssetRepository) DeleteForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Asset, error) {
	asset, err := r.FindByIDForUser(ec, id, userID)
	if err != nil {
		return model.Asset{}, err
	}

	err = r.database(ec).WithContext(ec.RequestContext()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM asset_vulnerabilities WHERE asset_id = ?", asset.ID).Error; err != nil {
			return err
		}
		return tx.Delete(&asset).Error
	})
	if err != nil {
		return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrDeleteFailed, err)
	}
	return asset, nil
}

func (r *AssetRepository) AssignVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, error) {
	asset, vulnerability, err := r.findAssetAndVulnerabilityForUser(ec, assetID, userID, vulnerabilityID)
	if err != nil {
		return model.Asset{}, err
	}

	for _, assigned := range asset.Vulnerabilities {
		if assigned.ID == vulnerability.ID {
			return model.Asset{}, baserepository.ErrDuplicateAssignment
		}
	}

	err = r.database(ec).WithContext(ec.RequestContext()).Model(&asset).Association("Vulnerabilities").Append(&vulnerability)
	if err != nil {
		databaseErr := utils.TranslateDatabaseError(err)
		if errors.Is(databaseErr, utils.ErrUniqueViolation) {
			return model.Asset{}, baserepository.ErrDuplicateAssignment
		}
		if errors.Is(databaseErr, utils.ErrForeignKeyViolation) {
			return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrInvalidReference, databaseErr)
		}
		return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrCreateFailed, databaseErr)
	}

	return r.FindByIDForUser(ec, assetID, userID)
}

func (r *AssetRepository) RemoveVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, error) {
	asset, vulnerability, err := r.findAssetAndVulnerabilityForUser(ec, assetID, userID, vulnerabilityID)
	if err != nil {
		return model.Asset{}, err
	}

	err = r.database(ec).WithContext(ec.RequestContext()).Model(&asset).Association("Vulnerabilities").Delete(&vulnerability)
	if err != nil {
		return model.Asset{}, fmt.Errorf("%w: %w", baserepository.ErrDeleteFailed, err)
	}

	return r.FindByIDForUser(ec, assetID, userID)
}

func (r *AssetRepository) findAssetAndVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, model.Vulnerability, error) {
	asset, err := r.FindByIDForUser(ec, assetID, userID)
	if err != nil {
		return model.Asset{}, model.Vulnerability{}, err
	}

	var vulnerability model.Vulnerability
	err = r.database(ec).WithContext(ec.RequestContext()).
		Where("user_id = ?", userID).
		First(&vulnerability, vulnerabilityID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Asset{}, model.Vulnerability{}, baserepository.ErrVulnerabilityNotFound
	}
	if err != nil {
		return model.Asset{}, model.Vulnerability{}, fmt.Errorf("%w: %w", baserepository.ErrReadFailed, err)
	}

	return asset, vulnerability, nil
}
