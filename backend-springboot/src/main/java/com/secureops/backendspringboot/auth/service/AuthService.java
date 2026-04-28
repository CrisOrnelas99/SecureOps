package com.secureops.backendspringboot.auth.service;

import com.secureops.backendspringboot.auth.dto.RegisterRequest;
import com.secureops.backendspringboot.auth.dto.LoginRequest;
import com.secureops.backendspringboot.auth.entity.User;
import com.secureops.backendspringboot.auth.repository.UserRepository;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    public AuthService(UserRepository userRepository, PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    public void register(RegisterRequest request) {

        if (userRepository.existsByUsername(request.getUsername()))
            throw new BadCredentialsException("Username is already taken.");

        if (userRepository.existsByEmail(request.getEmail()))
            throw new BadCredentialsException("Email is already in use.");

        User user = new User();
        user.setUsername(request.getUsername());
        user.setEmail(request.getEmail());
        user.setPasswordHash(passwordEncoder.encode(request.getPassword()));

        userRepository.save(user);
    }


    public void login(LoginRequest request) {
        User user = userRepository.findByUsername(request.getUserOrEmail())
                .or(() -> userRepository.findByEmail(request.getUserOrEmail()))
                .orElseThrow(() -> new BadCredentialsException("Invalid credentials."));

        if (!passwordEncoder.matches(request.getPassword(), user.getPasswordHash())) {
            throw new BadCredentialsException("Invalid credentials.");
        }
    }


}
