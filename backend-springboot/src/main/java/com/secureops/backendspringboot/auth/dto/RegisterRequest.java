package com.secureops.backendspringboot.auth.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public class RegisterRequest {

    @NotBlank
    @Size(min = 3, max = 20)
    private String username;
    public String getUsername() { return username; }
    public void setUsername(String username){
        this.username = username;
    }

    @NotBlank
    @Email
    private String email;
    public String getEmail() { return email; }
    public void setEmail(String email) {
        this.email = email;
    }

    @NotBlank
    @Size(min = 8, max = 100)
    private String password;
    public String getPassword(){ return password; }
    public void setPassword(String password) {
        this.password = password;
    }

}

//DTO - Data Transfer Object: class used to carry data between parts of the app
