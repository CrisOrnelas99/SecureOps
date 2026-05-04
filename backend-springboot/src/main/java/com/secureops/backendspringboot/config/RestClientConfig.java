package com.secureops.backendspringboot.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestClient;   //springs built in HTTP client for calling another API

@Configuration  //tells spring this class provides app configuration
public class RestClientConfig {

    @Bean   //tells spring to create one shared RestClient object
    public RestClient restClient() {

        return RestClient.builder().build();
    }
}