package model

type RiskCalculationRequest struct {
	AssetID                 int64  `json:"assetId"`
	Criticality             string `json:"criticality"`
	CriticalVulnerabilities int    `json:"criticalVulnerabilities"`
	HighVulnerabilities     int    `json:"highVulnerabilities"`
	MediumVulnerabilities   int    `json:"mediumVulnerabilities"`
	LowVulnerabilities      int    `json:"lowVulnerabilities"`
}

type RiskCalculationResponse struct {
	AssetID   int64  `json:"assetId"`
	RiskScore int    `json:"riskScore"`
	RiskLevel string `json:"riskLevel"`
}
