package dto

import (
	"strings"

	"secureops/backend-go/api/model"
)

type AssetRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	IPAddress       string  `json:"ipAddress"`
	OperatingSystem *string `json:"operatingSystem"`
	Owner           string  `json:"owner"`
	Criticality     string  `json:"criticality"`
}

func (r AssetRequest) ToDataModel() model.Asset {
	operatingSystem := trimOptionalString(r.OperatingSystem)

	return model.Asset{
		Name:            strings.TrimSpace(r.Name),
		Type:            strings.TrimSpace(r.Type),
		IPAddress:       strings.TrimSpace(r.IPAddress),
		OperatingSystem: operatingSystem,
		Owner:           strings.TrimSpace(r.Owner),
		Criticality:     strings.TrimSpace(r.Criticality),
		RiskScore:       0,
		RiskLevel:       "Low",
	}
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
