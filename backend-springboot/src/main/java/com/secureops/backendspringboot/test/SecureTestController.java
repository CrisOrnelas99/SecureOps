package com.secureops.backendspringboot.test;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class SecureTestController {

    @GetMapping("/api/test/secure")
    public String secureTest() {
        return "Secure endpoint reached";
    }
}