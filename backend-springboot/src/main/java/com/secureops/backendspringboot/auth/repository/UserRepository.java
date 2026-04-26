//Repository : the apps database helper for one entity

package com.secureops.backendspringboot.auth.repository;

import com.secureops.backendspringboot.auth.entity.User;
import org.springframework.data.jpa.repository.JpaRepository;   //gives built in database methods like save, findbyId, findAll, delete

import java.util.Optional;

public interface UserRepository extends JpaRepository<User, Long> {

    Optional<User> findByUsername(String username);

    Optional<User> findByEmail(String email);

    boolean existsByUsername(String username);  //checks whether a username already exists

    boolean existsByEmail(String email);    //checks whether an email already exists
}