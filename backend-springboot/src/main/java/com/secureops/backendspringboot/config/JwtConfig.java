package com.secureops.backendspringboot.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

@Configuration  //lets spring manage
public class JwtConfig {

    @Value("${jwt.secret}")    //pulls from application.properties
    private String secret;
    public String getSecret(){ return secret; }

    @Value("${jwt.expiration}")
    private long expiration;
    public long getExpiration(){ return expiration; }


}
