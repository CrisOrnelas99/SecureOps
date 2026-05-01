//what the client is asking to create/update

package com.secureops.backendspringboot.assets.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;

public class AssetRequest {

    @NotBlank
    private String name;
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }

    @NotBlank
    private String type;
    public String getType() { return type; }
    public void setType(String type) { this.type = type; }

    @NotBlank
    @Pattern(
            regexp = "^(?:25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(?:\\.(?:25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}$",
            message = "Invalid IPv4 address format."
    )
    private String ipAddress;
    public String getIpAddress() { return ipAddress; }
    public void setIpAddress(String ipAddress) { this.ipAddress = ipAddress; }

    private String operatingSystem;
    public String getOperatingSystem() { return operatingSystem; }
    public void setOperatingSystem(String operatingSystem) { this.operatingSystem = operatingSystem; }

    @NotBlank
    private String owner;
    public String getOwner() { return owner; }
    public void setOwner(String owner) { this.owner = owner; }

    @NotBlank
    private String criticality;
    public String getCriticality() { return criticality; }
    public void setCriticality(String criticality) { this.criticality = criticality; }

}
