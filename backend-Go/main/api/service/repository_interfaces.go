package service

import (
	"context"

	"secureops/backend-go/api/model"
)

type UserRepository interface {
	// ExistsByUsername reports whether a username is already stored.
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail reports whether an email address is already stored.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Save persists a new user record.
	Save(ctx context.Context, user model.User) error

	// FindByUsernameOrEmail returns the user matching a username or email.
	FindByUsernameOrEmail(ctx context.Context, userOrEmail string) (model.User, error)
}

type AssetRepository interface {
	// FindAll returns all stored assets ordered by ID.
	FindAll(ctx context.Context) ([]model.Asset, error)

	// FindByID returns one asset and its assigned vulnerabilities.
	FindByID(ctx context.Context, id int64) (model.Asset, error)

	// Save persists a new asset record.
	Save(ctx context.Context, request model.AssetRequest) (model.Asset, error)

	// Update changes an existing asset record.
	Update(ctx context.Context, id int64, request model.AssetRequest) (model.Asset, error)

	// Delete removes an asset record.
	Delete(ctx context.Context, id int64) (model.Asset, error)

	// AssignVulnerability links a vulnerability to an asset.
	AssignVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error)

	// RemoveVulnerability unlinks a vulnerability from an asset.
	RemoveVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error)

	// FindVulnerabilities returns the vulnerabilities assigned to an asset.
	FindVulnerabilities(ctx context.Context, assetID int64) ([]model.Vulnerability, error)

	// PersistRiskResult stores the calculated risk score and risk level for an asset.
	PersistRiskResult(ctx context.Context, id int64, riskResponse model.RiskCalculationResponse) (model.Asset, error)

	// EnsureAssetAndVulnerability confirms both records exist before assignment logic runs.
	EnsureAssetAndVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) error
}

type VulnerabilityRepository interface {
	// FindAll returns all stored vulnerabilities ordered by ID.
	FindAll(ctx context.Context) ([]model.Vulnerability, error)

	// FindByID returns one vulnerability by ID.
	FindByID(ctx context.Context, id int64) (model.Vulnerability, error)

	// Save persists a new vulnerability record.
	Save(ctx context.Context, request model.VulnerabilityRequest) (model.Vulnerability, error)

	// Update changes an existing vulnerability record.
	Update(ctx context.Context, id int64, request model.VulnerabilityRequest) (model.Vulnerability, error)

	// Delete removes a vulnerability record.
	Delete(ctx context.Context, id int64) (model.Vulnerability, error)
}
