package config

import "secureops/backend-go/api/service"

func RestClientConfig(cfg Config) *service.RestClient {
	return service.NewRestClient(cfg.RiskServiceURL, cfg.RiskServiceTimeout)
}
