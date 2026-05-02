package com.secureops.backendspringboot.assets.service;

import com.secureops.backendspringboot.assets.repository.AssetRepository;
import com.secureops.backendspringboot.assets.entity.Asset;
import com.secureops.backendspringboot.assets.dto.AssetRequest;
import org.springframework.stereotype.Service;
import java.util.List;

import com.secureops.backendspringboot.vulnerabilities.entity.Vulnerability;
import com.secureops.backendspringboot.vulnerabilities.repository.VulnerabilityRepository;

@Service
public class AssetService {

    private final AssetRepository assetRepository;
    private final VulnerabilityRepository vulnerabilityRepository;

    public AssetService(AssetRepository assetRepository, VulnerabilityRepository vulnerabilityRepository) {
        this.assetRepository = assetRepository;
        this.vulnerabilityRepository = vulnerabilityRepository;
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

    public Asset removeVulnerability(Long assetId, Long vulnerabilityId) {
        Asset asset = assetRepository.findById(assetId)
                .orElseThrow(() -> new IllegalArgumentException("Asset not found."));

        Vulnerability vulnerability = vulnerabilityRepository.findById(vulnerabilityId)
                .orElseThrow(() -> new IllegalArgumentException("Vulnerability not found."));

        asset.getVulnerabilities().remove(vulnerability);

        return assetRepository.save(asset);
    }


}
