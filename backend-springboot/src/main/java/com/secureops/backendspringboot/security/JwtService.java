package com.secureops.backendspringboot.security;

import com.secureops.backendspringboot.config.JwtConfig;
import io.jsonwebtoken.Claims;  //represent the data stored inside a JWT
import io.jsonwebtoken.Jwts;    //main JWT builder and parser utility
import io.jsonwebtoken.security.Keys;   //safely creates a signing key for JWT signatures
import org.springframework.stereotype.Service;  //marks this class as a Spring managed service

import javax.crypto.SecretKey;  //Type used for the HMAC signing key
import java.nio.charset.StandardCharsets;   //ensures the secret string is converted to bytes consistently
import java.util.Date;


@Service
public class JwtService {

    private final JwtConfig jwtConfig;

    public JwtService(JwtConfig jwtConfig) { this.jwtConfig = jwtConfig; }

    private SecretKey getSigningKey() { //builds secret key used to sign and verify JWT tokens
        return Keys.hmacShaKeyFor(jwtConfig.getSecret().getBytes(StandardCharsets.UTF_8));
    }

    // creates a signed JWT for the given username with expiration date
    public String generateToken(String username) {

        Date now = new Date();
        Date expiration = new Date(now.getTime() + jwtConfig.getExpiration());

        return Jwts.builder()
                .subject(username)
                .issuedAt(now)
                .expiration(expiration)
                .signWith(getSigningKey())
                .compact();
    }

    //reads the username stored in the token's subject field
    public String extractUsername(String token){
        return extractAllClaims(token).getSubject();
    }

    //confirms the token belongs to the expected user and is not expired
    public boolean isTokenValid(String token, String username){
        String extractedUsername = extractUsername(token);
        return extractedUsername.equals(username) && !isTokenExpired(token);
    }

    public boolean isTokenExpired(String token) {
        return extractAllClaims(token).getExpiration().before(new Date());
    }

    //parses the token, verifies the signature, and returns its stored claims
    private Claims extractAllClaims(String token) {
        return Jwts.parser()
                .verifyWith(getSigningKey())
                .build()
                .parseSignedClaims(token)
                .getPayload();
    }


}
