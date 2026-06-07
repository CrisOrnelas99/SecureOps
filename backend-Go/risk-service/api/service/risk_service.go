package service

import "risk-service-go/api/model"

const MaxVulnerabilityCount = 100000

type RiskService struct{}

func NewRiskService() *RiskService {
	return &RiskService{}
}

func (s *RiskService) Calculate(request model.RiskRequest) (model.RiskResponse, error) {
	if err := validate(request); err != nil {
		return model.RiskResponse{}, err
	}

	riskScore := (int64(request.CriticalVulnerabilities) * 25) +
		(int64(request.HighVulnerabilities) * 15) +
		(int64(request.MediumVulnerabilities) * 8) +
		(int64(request.LowVulnerabilities) * 3)

	switch request.Criticality {
	case "Medium":
		riskScore += 10
	case "High":
		riskScore += 20
	case "Critical":
		riskScore += 30
	}

	if riskScore > 100 {
		riskScore = 100
	}

	riskLevel := "Low"
	if riskScore >= 76 {
		riskLevel = "Critical"
	} else if riskScore >= 51 {
		riskLevel = "High"
	} else if riskScore >= 26 {
		riskLevel = "Medium"
	}

	return model.RiskResponse{
		AssetID:   request.AssetID,
		RiskScore: int(riskScore),
		RiskLevel: riskLevel,
	}, nil
}

func validate(request model.RiskRequest) error {
	if request.AssetID <= 0 {
		return invalidAssetID()
	}

	if request.Criticality != "Low" && request.Criticality != "Medium" && request.Criticality != "High" && request.Criticality != "Critical" {
		return invalidCriticality()
	}

	if request.CriticalVulnerabilities < 0 || request.HighVulnerabilities < 0 ||
		request.MediumVulnerabilities < 0 || request.LowVulnerabilities < 0 {
		return negativeVulnerabilities()
	}

	if request.CriticalVulnerabilities > MaxVulnerabilityCount || request.HighVulnerabilities > MaxVulnerabilityCount ||
		request.MediumVulnerabilities > MaxVulnerabilityCount || request.LowVulnerabilities > MaxVulnerabilityCount {
		return vulnerabilityLimitExceeded()
	}

	return nil
}
