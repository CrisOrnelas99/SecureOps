package com.secureops.backendspringboot.assets.repository;

import com.secureops.backendspringboot.assets.entity.Asset;
import org.springframework.data.jpa.repository.JpaRepository;

public interface AssetRepository extends JpaRepository<Asset, Long> {


}