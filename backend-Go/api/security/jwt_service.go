package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
	audience   string
}

type AccessClaims struct {
	Scope    string `json:"scope"`
	TokenUse string `json:"tokenUse"`
	jwt.RegisteredClaims
}

const (
	accessScope = "api"
	tokenUse    = "access"
)

func NewJWTManager(secret string, expiration time.Duration, issuer string, audience string) *JWTManager {
	return &JWTManager{secret: []byte(secret), expiration: expiration, issuer: issuer, audience: audience}
}

func (s *JWTManager) GenerateToken(username string) (string, error) {
	if len(s.secret) == 0 {
		return "", ErrMissingSecret
	}

	now := time.Now()
	claims := AccessClaims{
		Scope:    accessScope,
		TokenUse: tokenUse,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTManager) ExtractUsername(tokenString string) (string, error) {
	if len(s.secret) == 0 {
		return "", ErrMissingSecret
	}

	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedSigningMethod
		}
		return s.secret, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuer(s.issuer), jwt.WithAudience(s.audience), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}
	if claims.Subject == "" {
		return "", ErrMissingSubject
	}
	if claims.Scope != accessScope {
		return "", ErrInvalidScope
	}
	if claims.TokenUse != tokenUse {
		return "", ErrInvalidTokenUse
	}
	return claims.Subject, nil
}
