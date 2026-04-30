package com.secureops.backendspringboot.waf.entity;

import jakarta.persistence.*;
import java.time.LocalDateTime;

@Entity
@Table(name = "waf_events")
public class WafEvent {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    public Long getId() { return id; }

    private String method;
    public String getMethod() { return method; }
    public void setMethod(String method) { this.method = method; }

    private String path;
    public String getPath() { return path; }
    public void setPath(String path) { this.path = path; }

    private String reason;
    public String getReason() { return reason; }
    public void setReason(String reason) { this.reason = reason; }

    private String sourceIp;
    public String getSourceIp() { return sourceIp; }
    public void setSourceIp(String sourceIp) { this.sourceIp = sourceIp; }

    private LocalDateTime createdAt;
    public LocalDateTime getCreatedAt() { return createdAt; }
    public void setCreatedAt (LocalDateTime createdAt) { this.createdAt = createdAt; }
}
