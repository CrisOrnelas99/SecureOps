//spring configuration class that defines how passwords should be hashed across the app

package com.secureops.backendspringboot.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

@Configuration  //tells Spring this class provides application beans
public class PasswordConfig {

    @Bean   //tells spring to create one shared PasswordEncoder object and register it in app context
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(); //app uses BCrypt for password hashing
    }
}

//Bean : object that spring creates, owns, and gives to other parts of your app when needed