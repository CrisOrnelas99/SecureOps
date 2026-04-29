package com.secureops.backendspringboot.security;

import jakarta.servlet.http.HttpServletRequest; // Incoming request
import jakarta.servlet.http.HttpServletResponse; // Outgoing response
import org.springframework.security.core.AuthenticationException; // Spring auth failure type
import org.springframework.security.web.AuthenticationEntryPoint; // Handles unauthenticated access
import org.springframework.stereotype.Component; // Registers this as a Spring bean

import java.io.IOException; // Handles response writing errors

@Component
public class JwtAuthenticationEntryPoint implements AuthenticationEntryPoint {

    @Override
    public void commence(
            HttpServletRequest request,
            HttpServletResponse response,
            AuthenticationException authException
    ) throws IOException {

        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.setContentType("application/json");
        response.getWriter().write("{\"error\":\"Unauthorized\"}");
    }
}
