package com.secureops.backendspringboot.risk.dto;

public class RiskCalculationRequest {

    private Long assetId;
    public Long getAssetId() { return assetId; }
    public void setAssetId(Long assetId) { this.assetId = assetId; }

    private String criticality;
    public String getCriticality() { return criticality; }
    public void setCriticality(String criticality) { this.criticality = criticality; }

    private int criticalVulnerabilities;
    public int getCriticalVulnerabilities() { return criticalVulnerabilities; }
    public void setCriticalVulnerabilities(int criticalVulnerabilities) { this.criticalVulnerabilities = criticalVulnerabilities; }

    private int highVulnerabilities;
    public int getHighVulnerabilities() { return highVulnerabilities; }
    public void setHighVulnerabilities(int highVulnerabilities) {  this.highVulnerabilities = highVulnerabilities; }

    private int mediumVulnerabilities;
    public int getMediumVulnerabilities() { return mediumVulnerabilities; }
    public void setMediumVulnerabilities(int mediumVulnerabilities) { this.mediumVulnerabilities = mediumVulnerabilities; }

    private int lowVulnerabilities;
    public int getLowVulnerabilities() { return lowVulnerabilities; }
    public void setLowVulnerabilities(int lowVulnerabilities) { this.lowVulnerabilities = lowVulnerabilities; }




}