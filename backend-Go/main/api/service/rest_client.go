package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"secureops/backend-go/api/model"
)

type RestClient struct {
	baseURL string
	http    *http.Client
}

func NewRestClient(baseURL string, timeout time.Duration) *RestClient {
	return &RestClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *RestClient) CalculateRisk(ctx context.Context, request model.RiskCalculationRequest) (model.RiskCalculationResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return model.RiskCalculationResponse{}, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/calculate-risk", bytes.NewReader(body))
	if err != nil {
		return model.RiskCalculationResponse{}, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := c.http.Do(httpRequest)
	if err != nil {
		return model.RiskCalculationResponse{}, errors.Join(ErrRemoteService, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 400 && httpResponse.StatusCode < 500 {
		return model.RiskCalculationResponse{}, ErrRemoteRejected
	}
	if httpResponse.StatusCode >= 500 {
		return model.RiskCalculationResponse{}, ErrRemoteService
	}

	var result model.RiskCalculationResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&result); err != nil {
		return model.RiskCalculationResponse{}, errors.Join(ErrRemoteService, fmt.Errorf("%w: %v", ErrInvalidRemoteResult, err))
	}
	if result.RiskLevel == "" {
		return model.RiskCalculationResponse{}, errors.Join(ErrRemoteService, ErrInvalidRemoteResult)
	}

	return result, nil
}
