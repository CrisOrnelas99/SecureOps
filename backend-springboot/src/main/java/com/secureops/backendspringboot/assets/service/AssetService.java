package com.secureops.backendspringboot.assets.service;

import com.secureops.backendspringboot.assets.repository.AssetRepository;
import com.secureops.backendspringboot.assets.entity.Asset;
import com.secureops.backendspringboot.assets.dto.AssetRequest;
import org.springframework.stereotype.Service;
import java.util.List;

@Service
public class AssetService {

    private final AssetRepository assetRepository;

    public AssetService(AssetRepository assetRepository) {
        this.assetRepository = assetRepository;
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


}
