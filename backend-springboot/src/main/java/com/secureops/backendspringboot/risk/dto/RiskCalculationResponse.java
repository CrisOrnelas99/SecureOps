//DTO made specificically for backend to backend communication with the Go risk service

package com.secureops.backendspringboot.risk.dto;

public class RiskCalculationResponse {

    private Long assetId;
    public  Long getAssetId() { return assetId; }
    public void setAssetId(Long assetId) { this.assetId = assetId; }

    private int riskScore;
    public int getRiskScore() { return riskScore; }
    public void setRiskScore(int riskScore) { this.riskScore = riskScore; }

    private String riskLevel;
    public String getRiskLevel() { return riskLevel; }
    public void setRiskLevel(String riskLevel) { this.riskLevel = riskLevel; }


}