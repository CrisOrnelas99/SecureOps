package com.secureops.backendspringboot.assets.service;

import com.secureops.backendspringboot.assets.entity.Asset;
import com.secureops.backendspringboot.assets.repository.AssetRepository;
import com.secureops.backendspringboot.risk.dto.RiskCalculationRequest;
import com.secureops.backendspringboot.risk.dto.RiskCalculationResponse;
import com.secureops.backendspringboot.vulnerabilities.entity.Vulnerability;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class AssetRiskService {

    private final AssetRepository assetRepository;

    public AssetRiskService(AssetRepository assetRepository) {
        this.assetRepository = assetRepository;
    }

    @Transactional(readOnly = true)
    public RiskCalculationRequest loadRiskCalculationRequest(Long id) {
        Asset asset = assetRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        return buildRiskCalculationRequest(asset);
    }

    @Transactional
    public Asset persistRiskResult(Long id, RiskCalculationResponse response) {
        Asset asset = assetRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        if (response.getRiskLevel() == null) {
            throw new IllegalStateException("Risk service returned no risk level for asset id " + id + ".");
        }

        int riskScore = response.getRiskScore();
        if (riskScore < Short.MIN_VALUE || riskScore > Short.MAX_VALUE) {
            throw new IllegalStateException("Risk service returned an out-of-range risk score for asset id " + id + ": " + riskScore);
        }

        asset.setRiskScore((short) response.getRiskScore());
        asset.setRiskLevel(response.getRiskLevel());

        return assetRepository.save(asset);
    }

    private RiskCalculationRequest buildRiskCalculationRequest(Asset asset) {
        int criticalCount = 0;
        int highCount = 0;
        int mediumCount = 0;
        int lowCount = 0;

        for (Vulnerability vulnerability : asset.getVulnerabilities()) {
            String severity = vulnerability.getSeverity();
            if (severity == null) {
                continue;
            }

            String normalizedSeverity = severity.trim();

            if ("Critical".equalsIgnoreCase(normalizedSeverity))
                criticalCount++;
            else if ("High".equalsIgnoreCase(normalizedSeverity))
                highCount++;
            else if ("Medium".equalsIgnoreCase(normalizedSeverity))
                mediumCount++;
            else if ("Low".equalsIgnoreCase(normalizedSeverity))
                lowCount++;
        }

        RiskCalculationRequest request = new RiskCalculationRequest();
        request.setAssetId(asset.getId());
        request.setCriticality(asset.getCriticality());
        request.setCriticalVulnerabilities(criticalCount);
        request.setHighVulnerabilities(highCount);
        request.setMediumVulnerabilities(mediumCount);
        request.setLowVulnerabilities(lowCount);

        return request;
    }
}
