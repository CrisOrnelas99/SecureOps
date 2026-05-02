//what the backend stores in the database

package com.secureops.backendspringboot.assets.entity;

import jakarta.persistence.*;
import java.time.LocalDateTime;

import com.secureops.backendspringboot.vulnerabilities.entity.Vulnerability;
import java.util.HashSet;   //creates an empty set to hold assigned vulnerabilities
import java.util.Set;   //stores multiple unique vulnerabilites on one asset

@Entity
@Table(name="assets")
public class Asset {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    public Long getId() { return id; }

    @Column(nullable = false)
    private String name;
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }

    @Column(nullable = false)
    private String type;
    public String getType() { return type; }
    public void setType(String type) { this.type = type; }

    @Column(nullable = false)
    private String ipAddress;
    public String getIpAddress() { return ipAddress; }
    public void setIpAddress(String ipAddress) { this.ipAddress = ipAddress; }

    @Column
    private String operatingSystem;
    public String getOperatingSystem() { return operatingSystem; }
    public void setOperatingSystem(String operatingSystem) { this.operatingSystem = operatingSystem; }

    @Column(nullable = false)
    private String owner;
    public String getOwner() { return owner; }
    public void setOwner(String owner) { this.owner = owner; }

    @Column(nullable = false)
    private String criticality;
    public String getCriticality() { return criticality; }
    public void setCriticality(String criticality) { this.criticality = criticality; }

    @Column(nullable = false)
    private Short riskScore;
    public Short getRiskScore() { return riskScore; }
    public void setRiskScore(Short riskScore) { this.riskScore = riskScore; }

    @Column(nullable = false)
    private String riskLevel;
    public String getRiskLevel() { return riskLevel; }
    public void setRiskLevel(String riskLevel) { this.riskLevel = riskLevel; }

    @ManyToMany //One asset can be linkned to many vulnerabilities, and each vulnerabilit can affect many assets
    @JoinTable( //Defines the join table that connects assets and vulnerabilities
            name = "asset_vulnerabilities",
            joinColumns = @JoinColumn(name = "asset_id"),
            inverseJoinColumns = @JoinColumn(name = "vulnerability_id")
    )
    private Set<Vulnerability> vulnerabilities = new HashSet<>();   //holds the vulnerabilities aassigned to the asset
    public Set<Vulnerability> getVulnerabilities() { return vulnerabilities; }
    public void setVulnerabilities(Set<Vulnerability> vulnerabilities) {
        this.vulnerabilities = vulnerabilities;
    }

    @Column(nullable = false)
    private LocalDateTime createdAt;
    @Column(nullable = false)
    private LocalDateTime updatedAt;
    public LocalDateTime getCreatedAt() { return createdAt; }
    public LocalDateTime getUpdatedAt() { return updatedAt; }
    @PrePersist //runs right before a new entity is saved to the database for the first time
    public void prePersist() {
        LocalDateTime now = LocalDateTime.now();
        this.createdAt = now;
        this.updatedAt = now;
    }
    @PreUpdate  //runs right before an existing entity is updated in the database
    public void preUpdate() {
        this.updatedAt = LocalDateTime.now();
    }


}
