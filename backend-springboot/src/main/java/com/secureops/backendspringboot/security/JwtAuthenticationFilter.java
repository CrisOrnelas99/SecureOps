
//checks incoming request for a JWT token, and if the token is valid, tell Spring security which user is making the request

package com.secureops.backendspringboot.security;

import com.secureops.backendspringboot.auth.repository.UserRepository;   //loads the user from the database
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken; //creates an authenticated user object for spring security
import org.springframework.security.core.context.SecurityContextHolder; //stores auth info for the current request
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;  //attaches request details like IP/session info
import org.springframework.stereotype.Component;    //registers this filter as a spring managed bean
import org.springframework.web.filter.OncePerRequestFilter; //runs this filter once per HTTP request

import jakarta.servlet.FilterChain; //passes the request to the next filter
import jakarta.servlet.ServletException;    //handles servlet-level filter errors
import jakarta.servlet.http.HttpServletRequest; //represents the incoming HTTP request
import jakarta.servlet.http.HttpServletResponse;    //Represents the outgoing Http response

import java.io.IOException; //handles input/output errors in the filter
import io.jsonwebtoken.JwtException; // Catches invalid or tampered JWT parsing errors

@Component
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    private final JwtService jwtService;
    private final UserRepository userRepository;

    public JwtAuthenticationFilter(JwtService jwtService, UserRepository userRepository) {
        this.jwtService = jwtService;
        this.userRepository = userRepository;
    }

    //checks each request for a Bearer token and starts the authentication process if one is present
    @Override
    protected void doFilterInternal(
            HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain
        )
        throws ServletException, IOException{

            final String authHeader = request.getHeader("Authorization");   //Reads the Authorization header

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {  //Skips if no Bearer token is present
                filterChain.doFilter(request, response);
                return;
            }

            final String jwt = authHeader.substring(7); //removes "Bearer " and keeps only the token
            final String username;

            try {
                username = jwtService.extractUsername(jwt);    //Reads the username stored inside the token
            } catch (JwtException | IllegalArgumentException exception) {
                filterChain.doFilter(request, response); // Treats bad tokens as unauthenticated instead of crashing the request
                return;
            }


            if (username != null && SecurityContextHolder.getContext().getAuthentication() == null) {   //Only continue if a username exists and no user is already authenticated

                var user = userRepository.findByUsername(username).orElse(null);    //Looks up the user from the database

                if (user != null && jwtService.isTokenValid(jwt, user.getUsername())) { // Confirms the user exists and the token is still valid
                    UsernamePasswordAuthenticationToken authToken =
                            new UsernamePasswordAuthenticationToken(
                                    user.getUsername(),
                                    null,
                                    java.util.Collections.emptyList()
                            );

                    authToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request)); // Adds request details to the auth object

                    SecurityContextHolder.getContext().setAuthentication(authToken); // Marks this request as authenticated
                }
            }
            filterChain.doFilter(request, response);    //continues the request through the rest of spring security on to the controller
    }
}
