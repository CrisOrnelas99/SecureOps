
package com.secureops.backendspringboot.assets.controller;

import com.secureops.backendspringboot.assets.entity.Asset;
import com.secureops.backendspringboot.assets.service.AssetService;
import com.secureops.backendspringboot.assets.dto.AssetRequest;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/api/assets")
public class AssetController {

    private final AssetService assetService;

    public AssetController(AssetService assetService) { this.assetService = assetService; }

    @GetMapping
    public ResponseEntity<List<Asset>> getAssets() {
        List<Asset> assets = assetService.getAllAssets();
        return ResponseEntity.status(HttpStatus.OK).body(assets);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Asset> getAsset(@PathVariable Long id) {
        Asset retrievedAsset = assetService.getAsset(id);
        return ResponseEntity.status(HttpStatus.OK).body(retrievedAsset);
    }

    @PostMapping
    public ResponseEntity<Asset> createAsset(@Valid @RequestBody AssetRequest request) {
        Asset createdAsset = assetService.createAsset(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(createdAsset);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Asset> updateAsset(@PathVariable Long id, @Valid @RequestBody AssetRequest request) {
        Asset updatedAsset = assetService.updateAsset(id, request);
        return ResponseEntity.status(HttpStatus.OK).body(updatedAsset);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Asset> deleteAsset(@PathVariable Long id) {
        Asset deletedAsset = assetService.deleteAsset(id);
        return ResponseEntity.status(HttpStatus.OK).body(deletedAsset);
    }

}
