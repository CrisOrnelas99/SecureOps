package service

import (
	"context"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"secureops/backend-go/api/model"
	"secureops/backend-go/api/security"
)

type AuthService struct {
	userRepository UserRepository
	jwtService     *security.JwtService
}

func NewAuthService(userRepository UserRepository, jwtService *security.JwtService) *AuthService {
	return &AuthService{userRepository: userRepository, jwtService: jwtService}
}

func (s *AuthService) Register(ctx context.Context, request model.RegisterRequest) error {
	if err := validateRegisterRequest(request); err != nil {
		return err
	}

	exists, err := s.userRepository.ExistsByUsername(ctx, request.Username)
	if err != nil {
		return err
	}
	if exists {
		return ErrConflict
	}

	exists, err = s.userRepository.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepository.Save(ctx, model.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(hash),
	})
}

func (s *AuthService) Login(ctx context.Context, request model.LoginRequest) (string, error) {
	if strings.TrimSpace(request.UserOrEmail) == "" || utf8.RuneCountInString(request.Password) < 8 || utf8.RuneCountInString(request.Password) > 100 {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByUsernameOrEmail(ctx, request.UserOrEmail)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		return "", ErrInvalidCredentials
	}

	return s.jwtService.GenerateToken(user.Username)
}

func validateRegisterRequest(request model.RegisterRequest) error {
	usernameLen := utf8.RuneCountInString(request.Username)
	passwordLen := utf8.RuneCountInString(request.Password)

	if usernameLen < 3 || usernameLen > 20 {
		return ErrInvalidRequestData
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return ErrInvalidRequestData
	}
	if passwordLen < 8 || passwordLen > 100 {
		return ErrInvalidRequestData
	}

	return nil
}
