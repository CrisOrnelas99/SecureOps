package com.secureops.backendspringboot.exception;

public class ClientServiceException extends RuntimeException {

    public ClientServiceException(String message) {
        super(message);
    }
}
