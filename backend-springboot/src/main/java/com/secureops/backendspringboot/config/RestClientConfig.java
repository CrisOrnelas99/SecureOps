package com.secureops.backendspringboot.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.web.client.RestClient;   //springs built in HTTP client for calling another API

@Configuration  //tells spring this class provides app configuration
public class RestClientConfig {

    @Value("${risk.service.connect-timeout-ms:5000}")
    private int connectTimeoutMs;

    @Value("${risk.service.read-timeout-ms:5000}")
    private int readTimeoutMs;

    @Bean   //tells spring to create one shared RestClient object
    public RestClient restClient() {
        SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
        requestFactory.setConnectTimeout(connectTimeoutMs);
        requestFactory.setReadTimeout(readTimeoutMs);

        return RestClient.builder()
                .requestFactory(requestFactory)
                .build();
    }
}
