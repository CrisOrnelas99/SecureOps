package com.secureops.backendspringboot.auth.controller;

import com.secureops.backendspringboot.auth.dto.RegisterRequest;
import com.secureops.backendspringboot.auth.dto.LoginRequest;
import com.secureops.backendspringboot.auth.dto.LoginResponse;
import com.secureops.backendspringboot.auth.service.AuthService;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/auth")
public class AuthController {

    private final AuthService authService;

    public AuthController(AuthService authService) {
        this.authService = authService;
    }

    @PostMapping("/register")
    public ResponseEntity<String> register( @Valid @RequestBody RegisterRequest request) {
        authService.register(request);
        return ResponseEntity.status(HttpStatus.CREATED).body("User registered.");
    }

    //Verifies the login request and returns a JWT token for future authenticated requests
    @PostMapping("/login")
    public ResponseEntity<LoginResponse> login(@Valid @RequestBody LoginRequest request ){
        String token = authService.login(request);
        return ResponseEntity.ok(new LoginResponse(token, "Bearer"));

    }
}
