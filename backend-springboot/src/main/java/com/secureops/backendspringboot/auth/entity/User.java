package com.secureops.backendspringboot.auth.entity;

import jakarta.persistence.*;

@Entity //tells jpa this class maps to a database table
@Table(name = "users")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    public Long getId(){ return id; }

    @Column(nullable = false, unique = true)
    private String username;
    public String getUsername() { return username; }

    public void setUsername(String username) {
        this.username = username;
    }

    @Column(nullable = false, unique = true)
    private String email;
    public String getEmail() { return email; }
    public void setEmail(String email) {
        this.email = email;
    }

    @Column(nullable = false)
    private String passwordHash;
    public String getPasswordHash() { return passwordHash; }
    public void setPasswordHash(String passwordHash){
        this.passwordHash = passwordHash;
    }


}