package com.secureops.backendspringboot.waf.repository;

import com.secureops.backendspringboot.waf.entity.WafEvent;
import org.springframework.data.jpa.repository.JpaRepository;

public interface WafEventRepository extends JpaRepository<WafEvent, Long> {



}