package service

import "errors"

var (
	ErrInvalidAssetID          = errors.New("invalid asset id")
	ErrInvalidCriticality      = errors.New("invalid criticality")
	ErrNegativeVulnerabilities = errors.New("negative vulnerability count")
	ErrVulnerabilityLimit      = errors.New("vulnerability count exceeds limit")
)

type ServiceError struct {
	Kind    error
	Message string
}

func (e ServiceError) Error() string {
	return e.Message
}

func (e ServiceError) Unwrap() error {
	return e.Kind
}

func invalidAssetID() error {
	return ServiceError{Kind: ErrInvalidAssetID, Message: "assetId must be greater than 0"}
}

func invalidCriticality() error {
	return ServiceError{Kind: ErrInvalidCriticality, Message: "criticality must be Low, Medium, High, or Critical"}
}

func negativeVulnerabilities() error {
	return ServiceError{Kind: ErrNegativeVulnerabilities, Message: "vulnerability counts cannot be negative"}
}

func vulnerabilityLimitExceeded() error {
	return ServiceError{Kind: ErrVulnerabilityLimit, Message: "vulnerability counts exceed the maximum allowed value"}
}
