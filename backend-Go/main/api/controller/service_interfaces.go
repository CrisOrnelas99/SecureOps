package controller

import (
	"context"

	"secureops/backend-go/api/model"
)

type AuthService interface {
	// Register validates and creates a new user account.
	Register(ctx context.Context, request model.RegisterRequest) error

	// Login validates credentials and returns a JWT for the authenticated user.
	Login(ctx context.Context, request model.LoginRequest) (string, error)
}

type AssetService interface {
	// GetAllAssets returns all assets for the asset list endpoint.
	GetAllAssets(ctx context.Context) ([]model.Asset, error)

	// GetAsset returns one asset for the asset detail endpoint.
	GetAsset(ctx context.Context, id int64) (model.Asset, error)

	// CreateAsset validates and creates a new asset.
	CreateAsset(ctx context.Context, request model.AssetRequest) (model.Asset, error)

	// UpdateAsset validates and updates an existing asset.
	UpdateAsset(ctx context.Context, id int64, request model.AssetRequest) (model.Asset, error)

	// DeleteAsset removes an existing asset.
	DeleteAsset(ctx context.Context, id int64) (model.Asset, error)

	// AssignVulnerability links a vulnerability to an asset.
	AssignVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error)

	// RemoveVulnerability unlinks a vulnerability from an asset.
	RemoveVulnerability(ctx context.Context, assetID int64, vulnerabilityID int64) (model.Asset, error)

	// CalculateRisk recalculates and persists the asset risk result.
	CalculateRisk(ctx context.Context, id int64) (model.Asset, error)
}

type VulnerabilityService interface {
	// GetAllVulnerabilities returns all vulnerabilities for the vulnerability list endpoint.
	GetAllVulnerabilities(ctx context.Context) ([]model.Vulnerability, error)

	// GetVulnerability returns one vulnerability for the vulnerability detail endpoint.
	GetVulnerability(ctx context.Context, id int64) (model.Vulnerability, error)

	// CreateVulnerability validates and creates a new vulnerability.
	CreateVulnerability(ctx context.Context, request model.VulnerabilityRequest) (model.Vulnerability, error)

	// UpdateVulnerability validates and updates an existing vulnerability.
	UpdateVulnerability(ctx context.Context, id int64, request model.VulnerabilityRequest) (model.Vulnerability, error)

	// DeleteVulnerability removes an existing vulnerability.
	DeleteVulnerability(ctx context.Context, id int64) (model.Vulnerability, error)
}
