package model

type AssetRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	IPAddress       string  `json:"ipAddress"`
	OperatingSystem *string `json:"operatingSystem"`
	Owner           string  `json:"owner"`
	Criticality     string  `json:"criticality"`
}
