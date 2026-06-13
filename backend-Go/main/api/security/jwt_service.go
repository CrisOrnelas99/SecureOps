package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTManager(secret string, expiration time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiration: expiration}
}

func (s *JWTManager) GenerateToken(username string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTManager) ExtractUsername(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedSigningMethod
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}
	if claims.Subject == "" {
		return "", ErrMissingSubject
	}
	return claims.Subject, nil
}
