package com.secureops.backendspringboot.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import org.springframework.web.cors.CorsConfigurationSource;

import java.util.List;

//Cors: browsers block frontend requests to a different origin unless the backend allows them

@Configuration  //tells spring this class contains application configuration
public class CorsConfig {

    @Bean   //tells spring to create and manage the object returned by this method
    public CorsConfigurationSource coreConfigurationSource(){   //returns the cors rules spring should use
        CorsConfiguration config = new CorsConfiguration();

        config.setAllowedOrigins(List.of("http://localhost:4200"));
        config.setAllowedMethods(List.of("GET", "POST", "PUT", "DELETE"));
        config.setAllowedHeaders(List.of("Authorization", "Content-Type")); //for JWT and JSON requests
        config.setAllowCredentials(true);   //alows cookies or auth-related credentials to be included in requests

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource(); //creates an object that can apply CORS rules to URL paths
        source.registerCorsConfiguration("/**", config);    //apply CORS config to all backend routes

        return source;
    }

}

