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
	"secureops/backend-go/api/repository"
)

func TranslateRepositoryError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrAssetNotFound), errors.Is(err, repository.ErrVulnerabilityNotFound):
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	case errors.Is(err, repository.ErrDuplicateAssignment), errors.Is(err, repository.ErrDuplicateData), errors.Is(err, repository.ErrInvalidReference):
		return fmt.Errorf("%w: %w", ErrConflict, err)
	case errors.Is(err, repository.ErrInvalidData):
		return fmt.Errorf("%w: %w", ErrInvalidRequestData, err)
	default:
		return err
	}
}

func ValidateAsset(asset model.Asset) error {
	if strings.TrimSpace(asset.Name) == "" ||
		strings.TrimSpace(asset.Type) == "" ||
		strings.TrimSpace(asset.Owner) == "" ||
		strings.TrimSpace(asset.Criticality) == "" {
		return ErrInvalidRequestData
	}

	if ip := net.ParseIP(asset.IPAddress); ip == nil || ip.To4() == nil {
		return ErrInvalidRequestData
	}

	return nil
}

func AuthenticatedUserID(ec *appcontext.GinContext) (int64, error) {
	if ec == nil || ec.UserID() <= 0 {
		return 0, ErrForbidden
	}
	return ec.UserID(), nil
}

func NormalizeRegisterRequest(request dto.RegisterRequest) dto.RegisterRequest {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	return request
}

func ValidateRegisterRequest(request dto.RegisterRequest) error {
	usernameLen := utf8.RuneCountInString(request.Username)
	passwordLen := utf8.RuneCountInString(request.Password)

	if usernameLen < 3 || usernameLen > 20 {
		return ErrInvalidRequestData
	}
	address, err := mail.ParseAddress(request.Email)
	if err != nil || address.Name != "" || address.Address != request.Email {
		return ErrInvalidRequestData
	}
	if passwordLen < 8 || passwordLen > 100 {
		return ErrInvalidRequestData
	}

	return nil
}

func ValidateVulnerability(vulnerability model.Vulnerability) error {
	if strings.TrimSpace(vulnerability.CVEID) == "" ||
		strings.TrimSpace(vulnerability.Title) == "" ||
		strings.TrimSpace(vulnerability.Description) == "" {
		return ErrInvalidRequestData
	}
	if vulnerability.Severity != "Low" && vulnerability.Severity != "Medium" && vulnerability.Severity != "High" && vulnerability.Severity != "Critical" {
		return ErrInvalidRequestData
	}
	if vulnerability.Status != "Open" && vulnerability.Status != "Fixed" && vulnerability.Status != "In Progress" {
		return ErrInvalidRequestData
	}
	return nil
}
