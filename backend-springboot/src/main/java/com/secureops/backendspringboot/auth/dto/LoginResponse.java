package com.secureops.backendspringboot.auth.dto;

public class LoginResponse {

    private String token;
    private String type;

    public LoginResponse(String token, String type) {
        this.token = token;
        this.type = type;
    }

    public String getToken() { return token; }

    public String getType() { return type; }
}