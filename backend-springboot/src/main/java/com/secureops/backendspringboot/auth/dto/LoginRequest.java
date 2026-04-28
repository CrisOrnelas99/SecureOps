package com.secureops.backendspringboot.auth.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public class LoginRequest{

    @NotBlank
    private String userOrEmail;
    public String getUserOrEmail() { return userOrEmail; }
    public void setUserOrEmail( String userOrEmail ) { this.userOrEmail = userOrEmail; }

    @NotBlank
    @Size(min=8, max=100)
    private String password;
    public String getPassword() { return password; }
    public void setPassword( String password ) { this.password = password; }

}