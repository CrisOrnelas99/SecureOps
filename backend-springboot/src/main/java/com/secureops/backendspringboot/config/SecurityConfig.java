package com.secureops.backendspringboot.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;

import com.secureops.backendspringboot.security.JwtAuthenticationFilter; // Runs JWT checks before Spring's built-in auth filter
import org.springframework.security.config.http.SessionCreationPolicy; // Makes authentication stateless for JWT APIs
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter; // Used to place the JWT filter in the correct order
import com.secureops.backendspringboot.security.JwtAuthenticationEntryPoint;

@Configuration
public class SecurityConfig {

    private final JwtAuthenticationFilter jwtAuthenticationFilter;
    private final JwtAuthenticationEntryPoint jwtAuthenticationEntryPoint;

    public SecurityConfig(
            JwtAuthenticationFilter jwtAuthenticationFilter,
            JwtAuthenticationEntryPoint jwtAuthenticationEntryPoint
    ){
        this.jwtAuthenticationFilter = jwtAuthenticationFilter;
        this.jwtAuthenticationEntryPoint = jwtAuthenticationEntryPoint;
    }

    //tells spring to create and manage the object returned by this method
    @Bean
    //defines the main security rules for incoming HTTP requests
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
                .csrf(csrf -> csrf.disable())   //turns off CSRF protection
                .cors(Customizer.withDefaults())    //tells spring security to use your cors config
                .exceptionHandling(exception -> exception.authenticationEntryPoint(jwtAuthenticationEntryPoint)) // Returns 401 when authentication is missing or invalid
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/api/health", "/api/auth/register", "/api/auth/login").permitAll() // Public endpoints stay open
                        .anyRequest().authenticated() // Everything else requires authentication
                )
                .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS)) // Prevents Spring from creating login sessions
                .addFilterBefore(jwtAuthenticationFilter, UsernamePasswordAuthenticationFilter.class); // Runs JWT auth before Spring's default auth filter

        return http.build();
    }
}

//HttpSecurity : the object you use to configure those rules
//SecurityFilterChain : final security setup spring will apply

