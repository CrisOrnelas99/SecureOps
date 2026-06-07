package repository

type RepositoryError struct {
	Message string
}

func (e RepositoryError) Error() string {
	return e.Message
}

var (
	ErrAssetNotFound         = &RepositoryError{Message: "asset not found"}
	ErrVulnerabilityNotFound = &RepositoryError{Message: "vulnerability not found"}
	ErrDuplicateAssignment   = &RepositoryError{Message: "duplicate asset vulnerability assignment"}
	ErrRiskScoreOutOfRange   = &RepositoryError{Message: "risk score out of range"}
)
