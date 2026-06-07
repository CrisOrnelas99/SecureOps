package service

import (
	"errors"

	"secureops/backend-go/api/repository"
)

func mapRepositoryError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrAssetNotFound), errors.Is(err, repository.ErrVulnerabilityNotFound):
		return ErrNotFound
	case errors.Is(err, repository.ErrDuplicateAssignment):
		return ErrConflict
	case errors.Is(err, repository.ErrRiskScoreOutOfRange):
		return ErrInvalidRequestData
	default:
		return err
	}
}
