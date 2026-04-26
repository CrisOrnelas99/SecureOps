package com.secureops.backendspringboot.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
public class SecurityConfig {

    //tells spring to create and manage the object returned by this method
    @Bean
    //defines the main security rules for incoming HTTP requests
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
                .csrf(csrf -> csrf.disable())   //turns off CSRF protection
                .cors(Customizer.withDefaults())    //tells spring security to use your cors config
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/api/health", "/api/auth/register").permitAll() //for test purposes
                        .anyRequest().authenticated()   //keeps the backedn secure by default
                )
                .httpBasic(Customizer.withDefaults());  //enables basic HTTP authentication for now

        return http.build();

    }
}

//HttpSecurity : the object you use to configure those rules
//SecurityFilterChain : final security setup spring will apply

