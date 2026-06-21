// Package service provides authentication application services.
package service

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"secureops/backend-go/api/config"
	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
	"secureops/backend-go/api/model"
	baserepository "secureops/backend-go/api/repository"
	"secureops/backend-go/api/security"
	baseservice "secureops/backend-go/api/service"
)

type authServiceImpl struct {
	jwtManager     *security.JWTManager
	userRepository baserepository.UserRepository
}

// NewAuthService creates an authentication service backed by the supplied dependencies.
func NewAuthService(jwtManager *security.JWTManager, userRepository baserepository.UserRepository) baseservice.AuthService {
	return &authServiceImpl{jwtManager: jwtManager, userRepository: userRepository}
}

// Register validates and creates a new user account.
func (s *authServiceImpl) Register(ec *appcontext.GinContext, request dto.RegisterRequest) error {
	request = baseservice.NormalizeRegisterRequest(request)
	if err := baseservice.ValidateRegisterRequest(request); err != nil {
		return err
	}

	exists, err := s.userRepository.ExistsByUsername(ec, request.Username)
	if err != nil {
		return baseservice.TranslateRepositoryError(err)
	}
	if exists {
		return baseservice.ErrConflict
	}

	exists, err = s.userRepository.ExistsByEmail(ec, request.Email)
	if err != nil {
		return baseservice.TranslateRepositoryError(err)
	}
	if exists {
		return baseservice.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), config.PasswordCost())
	if err != nil {
		return err
	}

	return baseservice.TranslateRepositoryError(s.userRepository.Save(ec, model.User{
		Username:     request.Username,
		Email:        request.Email,
		Role:         model.RoleUser,
		PasswordHash: string(hash),
	}))
}

// Login validates credentials and returns a signed access token.
func (s *authServiceImpl) Login(ec *appcontext.GinContext, request dto.LoginRequest) (dto.LoginResponse, error) {
	request.UserOrEmail = strings.TrimSpace(request.UserOrEmail)
	if strings.Contains(request.UserOrEmail, "@") {
		request.UserOrEmail = strings.ToLower(request.UserOrEmail)
	}
	if request.UserOrEmail == "" || utf8.RuneCountInString(request.Password) < 8 || utf8.RuneCountInString(request.Password) > 100 {
		return dto.LoginResponse{}, baseservice.ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByUsernameOrEmail(ec, request.UserOrEmail)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.LoginResponse{}, baseservice.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, baseservice.TranslateRepositoryError(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		return dto.LoginResponse{}, baseservice.ErrInvalidCredentials
	}

	if s.jwtManager == nil {
		return dto.LoginResponse{}, fmt.Errorf("missing jwt manager")
	}

	token, err := s.jwtManager.GenerateToken(user.Username)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}
