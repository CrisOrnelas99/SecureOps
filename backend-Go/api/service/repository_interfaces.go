package service

import (
	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/model"
)

type UserRepository interface {
	// ExistsByUsername reports whether a username is already stored.
	ExistsByUsername(ec *appcontext.GinContext, username string) (bool, error)

	// ExistsByEmail reports whether an email address is already stored.
	ExistsByEmail(ec *appcontext.GinContext, email string) (bool, error)

	// Save persists a new user record.
	Save(ec *appcontext.GinContext, user model.User) error

	// FindByUsernameOrEmail returns the user matching a username or email.
	FindByUsernameOrEmail(ec *appcontext.GinContext, userOrEmail string) (model.User, error)
}

type AssetRepository interface {
	// FindAllByUser returns all assets owned by a user ordered by ID.
	FindAllByUser(ec *appcontext.GinContext, userID int64) ([]model.Asset, error)

	// FindByIDForUser returns one owned asset and its assigned vulnerabilities.
	FindByIDForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Asset, error)

	// Save persists a new asset record.
	Save(ec *appcontext.GinContext, asset model.Asset) (model.Asset, error)

	// UpdateForUser changes an existing owned asset record.
	UpdateForUser(ec *appcontext.GinContext, id int64, userID int64, asset model.Asset) (model.Asset, error)

	// DeleteForUser removes an owned asset record.
	DeleteForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Asset, error)

	// AssignVulnerabilityForUser links a vulnerability to an owned asset.
	AssignVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, error)

	// RemoveVulnerabilityForUser unlinks a vulnerability from an owned asset.
	RemoveVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, error)
}

type VulnerabilityRepository interface {
	// FindAllByUser returns all vulnerabilities owned by a user ordered by ID.
	FindAllByUser(ec *appcontext.GinContext, userID int64) ([]model.Vulnerability, error)

	// FindByIDForUser returns one owned vulnerability by ID.
	FindByIDForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Vulnerability, error)

	// ExistsByCVEIDForUser reports whether a user already owns the CVE ID.
	ExistsByCVEIDForUser(ec *appcontext.GinContext, cveID string, userID int64) (bool, error)

	// ExistsByCVEIDExcludingIDForUser reports whether another owned vulnerability has the CVE ID.
	ExistsByCVEIDExcludingIDForUser(ec *appcontext.GinContext, cveID string, id int64, userID int64) (bool, error)

	// Save persists a new vulnerability record.
	Save(ec *appcontext.GinContext, vulnerability model.Vulnerability) (model.Vulnerability, error)

	// UpdateForUser changes an existing owned vulnerability record.
	UpdateForUser(ec *appcontext.GinContext, id int64, userID int64, vulnerability model.Vulnerability) (model.Vulnerability, error)

	// DeleteForUser removes an owned vulnerability record.
	DeleteForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Vulnerability, error)
}
