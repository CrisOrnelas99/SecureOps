// Package service provides validation, context, and repository error helpers for application services.
package service

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"strings"
	"unicode/utf8"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
	"secureops/backend-go/api/model"
	baserepository "secureops/backend-go/api/repository"
)

// TranslateRepositoryError maps repository errors to service-layer sentinels.
func TranslateRepositoryError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, baserepository.ErrAssetNotFound), errors.Is(err, baserepository.ErrVulnerabilityNotFound):
		return ErrNotFound
	case errors.Is(err, baserepository.ErrDuplicateData), errors.Is(err, baserepository.ErrDuplicateAssignment):
		return ErrConflict
	case errors.Is(err, baserepository.ErrInvalidData), errors.Is(err, baserepository.ErrInvalidReference):
		return ErrInvalidRequestData
	default:
		return err
	}
}

// ValidateAsset validates the fields required to create or update an asset.
func ValidateAsset(asset model.Asset) error {
	if strings.TrimSpace(asset.Name) == "" || strings.TrimSpace(asset.Type) == "" || strings.TrimSpace(asset.IPAddress) == "" || strings.TrimSpace(asset.Owner) == "" || strings.TrimSpace(asset.Criticality) == "" {
		return ErrInvalidRequestData
	}
	if net.ParseIP(strings.TrimSpace(asset.IPAddress)) == nil {
		return ErrInvalidRequestData
	}
	return nil
}

// AuthenticatedUserID returns the authenticated user ID from the request context.
func AuthenticatedUserID(ec *appcontext.GinContext) (int64, error) {
	if ec == nil {
		return 0, ErrForbidden
	}

	userID := ec.UserID()
	if userID <= 0 {
		return 0, ErrForbidden
	}

	return userID, nil
}

// NormalizeRegisterRequest trims and normalizes registration input.
func NormalizeRegisterRequest(request dto.RegisterRequest) dto.RegisterRequest {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Password = strings.TrimSpace(request.Password)
	return request
}

// ValidateRegisterRequest validates the fields required to create an account.
func ValidateRegisterRequest(request dto.RegisterRequest) error {
	if strings.TrimSpace(request.Username) == "" || utf8.RuneCountInString(request.Username) < 3 || utf8.RuneCountInString(request.Username) > 50 {
		return ErrInvalidRequestData
	}
	if strings.TrimSpace(request.Password) == "" || utf8.RuneCountInString(request.Password) < 8 || utf8.RuneCountInString(request.Password) > 100 {
		return ErrInvalidRequestData
	}
	if strings.TrimSpace(request.Email) == "" {
		return ErrInvalidRequestData
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return fmt.Errorf("%w: invalid email", ErrInvalidRequestData)
	}
	return nil
}

// ValidateVulnerability validates the fields required to create or update a vulnerability.
func ValidateVulnerability(vulnerability model.Vulnerability) error {
	if strings.TrimSpace(vulnerability.CVEID) == "" || strings.TrimSpace(vulnerability.Title) == "" || strings.TrimSpace(vulnerability.Severity) == "" || strings.TrimSpace(vulnerability.Description) == "" || strings.TrimSpace(vulnerability.Status) == "" {
		return ErrInvalidRequestData
	}
	return nil
}
