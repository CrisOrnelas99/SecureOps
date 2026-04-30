package com.secureops.backendspringboot.security;

import com.secureops.backendspringboot.waf.entity.WafEvent;
import com.secureops.backendspringboot.waf.repository.WafEventRepository;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.time.LocalDateTime;

@Component
public class WafFilter extends OncePerRequestFilter {

    private static final org.slf4j.Logger logger = org.slf4j.LoggerFactory.getLogger(WafFilter.class); // Creates a logger for this filter
    private final WafEventRepository wafEventRepository;

    public WafFilter(WafEventRepository wafEventRepository){
        this.wafEventRepository = wafEventRepository;
    }

    @Override
    protected void doFilterInternal(
            HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain
    )
        throws ServletException, IOException {
            String requestUri = request.getRequestURI(); //gets the request path
            String queryString = request.getQueryString();  //gets the raw query part after ?, if present
            String dataToInspect = (requestUri + " " + (queryString != null ? queryString : "")).toLowerCase(); //
            //Combines path and query text into one lowercase string for simple matching
            String reason = null;

            boolean hasPathTraversal = dataToInspect.contains("../"); // Checks for path traversal attempts
            boolean hasXssPattern =
                    dataToInspect.contains("<script") ||
                    dataToInspect.contains("%3cscript"); // Checks for basic XSS patterns, including URL-encoded input
            boolean hasSqlInjectionPattern =
                    dataToInspect.contains("' or ") ||
                    dataToInspect.contains("%27%20or%20") ||
                    dataToInspect.contains("union select") ||
                    dataToInspect.contains("drop table"); // Checks for a few obvious SQL injection-like strings

            if (hasPathTraversal || hasXssPattern || hasSqlInjectionPattern) {
                logger.warn("Blocked suspicious request: method={}, path={}", request.getMethod(),
                        request.getRequestURI()); // Logs the HTTP method and path without logging the full payload

                if (hasPathTraversal)
                    reason = "PATH_TRAVERSAL";
                else if (hasXssPattern)
                    reason = "XSS_PATTERN";
                else if (hasSqlInjectionPattern)
                    reason = "SQLI_PATTERN";

                WafEvent wafEvent = new WafEvent();
                wafEvent.setMethod(request.getMethod());
                wafEvent.setPath(request.getRequestURI());
                wafEvent.setReason(reason);
                wafEvent.setSourceIp(request.getRemoteAddr());
                wafEvent.setCreatedAt(LocalDateTime.now());

                wafEventRepository.save(wafEvent);

                response.setStatus(HttpServletResponse.SC_FORBIDDEN); // Returns HTTP 403 Forbidden
                response.setContentType("text/plain"); // Sends a simple plain-text response
                response.getWriter().write("Request blocked"); // Keeps the response generic and safe
                return;
            }

            filterChain.doFilter(request, response); //lets the request continue for now

    }

}

//@Override : this methos is replacing a method that already exists in the parent class
//ServletException : general web/filter exception used by the Java servlet system
//IOException : an input/output operation failed
