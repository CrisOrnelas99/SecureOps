package com.secureops.backendspringboot.assets.service;

import com.secureops.backendspringboot.assets.repository.AssetRepository;
import com.secureops.backendspringboot.assets.entity.Asset;
import com.secureops.backendspringboot.assets.dto.AssetRequest;
import org.springframework.stereotype.Service;
import java.util.List;
import org.springframework.transaction.annotation.Transactional;

import com.secureops.backendspringboot.vulnerabilities.entity.Vulnerability;
import com.secureops.backendspringboot.vulnerabilities.repository.VulnerabilityRepository;

import com.secureops.backendspringboot.risk.dto.RiskCalculationRequest;
import com.secureops.backendspringboot.risk.dto.RiskCalculationResponse;
import com.secureops.backendspringboot.exception.ClientServiceException;
import com.secureops.backendspringboot.exception.RemoteServiceException;
import org.springframework.web.client.RestClient;   //lets your spring boot backend make HTTP requests to another service
import org.springframework.web.client.RestClientException;
import org.springframework.util.StreamUtils;

import java.nio.charset.StandardCharsets;

@Service
public class AssetService {

    private final AssetRepository assetRepository;
    private final VulnerabilityRepository vulnerabilityRepository;
    private final RestClient restClient;
    private final AssetRiskService assetRiskService;

    public AssetService(
            AssetRepository assetRepository,
            VulnerabilityRepository vulnerabilityRepository,
            RestClient restClient,
            AssetRiskService assetRiskService
    ) {
        this.assetRepository = assetRepository;
        this.vulnerabilityRepository = vulnerabilityRepository;
        this.restClient = restClient;
        this.assetRiskService = assetRiskService;
    }

    public List<Asset> getAllAssets() {

        return assetRepository.findAll();
    }

    public Asset getAsset(Long id) {

        Asset asset = assetRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        return asset;
    }

    public Asset createAsset(AssetRequest request){

        Asset asset = new Asset();

        asset.setName(request.getName());
        asset.setType(request.getType());
        asset.setIpAddress(request.getIpAddress());
        asset.setOperatingSystem(request.getOperatingSystem());
        asset.setOwner(request.getOwner());
        asset.setRiskScore((short) 0);
        asset.setRiskLevel("Low");
        asset.setCriticality(request.getCriticality());

        return assetRepository.save(asset);
    }

    public Asset updateAsset(Long id, AssetRequest request){

        Asset asset = assetRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        asset.setName(request.getName());
        asset.setType(request.getType());
        asset.setIpAddress(request.getIpAddress());
        asset.setOperatingSystem(request.getOperatingSystem());
        asset.setOwner(request.getOwner());
        asset.setCriticality(request.getCriticality());

        return assetRepository.save(asset);
    }

    public Asset deleteAsset(Long id) {
        Asset asset = assetRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        assetRepository.delete(asset);
        return asset;
    }

    @Transactional  //wraps the entire method in a single transaction, making the check-then-insert atomic and eliminating session dependency
    public Asset assignVulnerability(Long assetId, Long vulnerabilityId) {
        Asset asset = assetRepository.findById(assetId)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        Vulnerability vulnerability = vulnerabilityRepository.findById(vulnerabilityId)
                .orElseThrow(() -> new IllegalArgumentException("Vulnerability not found."));

        if (asset.getVulnerabilities().contains(vulnerability)) {
            throw new IllegalArgumentException("Vulnerability is already assigned to this asset.");
        }

        asset.getVulnerabilities().add(vulnerability);

        return assetRepository.save(asset);
    }
    @Transactional
    public Asset removeVulnerability(Long assetId, Long vulnerabilityId) {
        Asset asset = assetRepository.findById(assetId)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        Vulnerability vulnerability = vulnerabilityRepository.findById(vulnerabilityId)
                .orElseThrow(() -> new IllegalArgumentException("Vulnerability not found."));

        asset.getVulnerabilities().remove(vulnerability);

        return assetRepository.save(asset);
    }

    private RiskCalculationRequest buildRiskCalculationRequest(Asset asset) {

        int criticalCount = 0;
        int highCount = 0;
        int mediumCount = 0;
        int lowCount = 0;

        for (Vulnerability vulnerability : asset.getVulnerabilities()) {    //loops through the assets assigned vulnerabilities

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

    public Asset calculateRisk(Long id) {
        RiskCalculationRequest request = assetRiskService.loadRiskCalculationRequest(id);

        RiskCalculationResponse response;
        try {
            response = restClient.post()
                    .uri("http://risk-service:8081/calculate-risk")
                    .body(request)
                    .retrieve()
                    .onStatus(
                            status -> status.is4xxClientError(),
                            (clientRequest, clientResponse) -> {
                                String responseBody = StreamUtils.copyToString(clientResponse.getBody(), StandardCharsets.UTF_8);
                                throw new ClientServiceException("Risk service rejected the request: " + responseBody);
                            }
                    )
                    .onStatus(
                            status -> status.is5xxServerError(),
                            (clientRequest, clientResponse) -> {
                                String responseBody = StreamUtils.copyToString(clientResponse.getBody(), StandardCharsets.UTF_8);
                                throw new RemoteServiceException("Risk service failed: " + responseBody);
                            }
                    )
                    .body(RiskCalculationResponse.class);
        } catch (RestClientException ex) {
            throw new RemoteServiceException("Failed to call risk service.", ex);
        }

        if (response == null) {
            throw new IllegalStateException("Risk service returned no response.");
        }

        return assetRiskService.persistRiskResult(id, response);
    }


}
