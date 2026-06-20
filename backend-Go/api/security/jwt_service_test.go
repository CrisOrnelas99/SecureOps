package security

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTManagerGenerateTokenAndExtractUsername(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, "issuer", "audience")

	token, err := manager.GenerateToken("analyst")
	if err != nil {
		t.Fatalf("expected token generation to succeed: %v", err)
	}
	if token == "" {
		t.Fatal("expected generated token to be non-empty")
	}

	username, err := manager.ExtractUsername(token)
	if err != nil {
		t.Fatalf("expected username extraction to succeed: %v", err)
	}
	if username != "analyst" {
		t.Fatalf("expected username analyst, got %q", username)
	}
}

func TestJWTManagerGenerateTokenIncludesExpectedClaims(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, "issuer", "audience")

	tokenString, err := manager.GenerateToken("analyst")
	if err != nil {
		t.Fatalf("expected token generation to succeed: %v", err)
	}

	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte("test-secret"), nil
	}, jwt.WithIssuer("issuer"), jwt.WithAudience("audience"), jwt.WithExpirationRequired())
	if err != nil {
		t.Fatalf("expected generated token to parse: %v", err)
	}
	if !token.Valid {
		t.Fatal("expected generated token to be valid")
	}
	if claims.Subject != "analyst" {
		t.Fatalf("expected subject analyst, got %q", claims.Subject)
	}
	if claims.Scope != accessScope {
		t.Fatalf("expected scope %q, got %q", accessScope, claims.Scope)
	}
	if claims.TokenUse != tokenUse {
		t.Fatalf("expected token use %q, got %q", tokenUse, claims.TokenUse)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expected expiration claim to be set")
	}
}

func TestJWTManagerRequiresSecret(t *testing.T) {
	manager := NewJWTManager("", time.Hour, "issuer", "audience")

	token, err := manager.GenerateToken("analyst")
	if !errors.Is(err, ErrMissingSecret) {
		t.Fatalf("expected ErrMissingSecret from GenerateToken, got token=%q err=%v", token, err)
	}

	username, err := manager.ExtractUsername("token")
	if !errors.Is(err, ErrMissingSecret) {
		t.Fatalf("expected ErrMissingSecret from ExtractUsername, got username=%q err=%v", username, err)
	}
}

func TestJWTManagerRejectsInvalidTokens(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, "issuer", "audience")

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "malformed token",
			token: "not-a-token",
		},
		{
			name: "wrong signing secret",
			token: signToken(t, "wrong-secret", AccessClaims{
				Scope:            accessScope,
				TokenUse:         tokenUse,
				RegisteredClaims: validRegisteredClaims("analyst"),
			}),
		},
		{
			name: "expired token",
			token: signToken(t, "test-secret", AccessClaims{
				Scope:    accessScope,
				TokenUse: tokenUse,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "analyst",
					Issuer:    "issuer",
					Audience:  jwt.ClaimStrings{"audience"},
					IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				},
			}),
		},
		{
			name: "wrong issuer",
			token: signToken(t, "test-secret", AccessClaims{
				Scope:    accessScope,
				TokenUse: tokenUse,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "analyst",
					Issuer:    "other-issuer",
					Audience:  jwt.ClaimStrings{"audience"},
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
			}),
		},
		{
			name: "wrong audience",
			token: signToken(t, "test-secret", AccessClaims{
				Scope:    accessScope,
				TokenUse: tokenUse,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "analyst",
					Issuer:    "issuer",
					Audience:  jwt.ClaimStrings{"other-audience"},
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
			}),
		},
		{
			name: "missing expiration",
			token: signToken(t, "test-secret", AccessClaims{
				Scope:    accessScope,
				TokenUse: tokenUse,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:  "analyst",
					Issuer:   "issuer",
					Audience: jwt.ClaimStrings{"audience"},
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := manager.ExtractUsername(tt.token)
			if !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("expected ErrInvalidToken, got username=%q err=%v", username, err)
			}
		})
	}
}

func TestJWTManagerRejectsInvalidApplicationClaims(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour, "issuer", "audience")

	tests := []struct {
		name      string
		claims    AccessClaims
		expectErr error
	}{
		{
			name: "missing subject",
			claims: AccessClaims{
				Scope:    accessScope,
				TokenUse: tokenUse,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "issuer",
					Audience:  jwt.ClaimStrings{"audience"},
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
			},
			expectErr: ErrMissingSubject,
		},
		{
			name: "invalid scope",
			claims: AccessClaims{
				Scope:            "admin",
				TokenUse:         tokenUse,
				RegisteredClaims: validRegisteredClaims("analyst"),
			},
			expectErr: ErrInvalidScope,
		},
		{
			name: "invalid token use",
			claims: AccessClaims{
				Scope:            accessScope,
				TokenUse:         "refresh",
				RegisteredClaims: validRegisteredClaims("analyst"),
			},
			expectErr: ErrInvalidTokenUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := signToken(t, "test-secret", tt.claims)

			username, err := manager.ExtractUsername(token)
			if !errors.Is(err, tt.expectErr) {
				t.Fatalf("expected %v, got username=%q err=%v", tt.expectErr, username, err)
			}
		})
	}
}

func TestSecurityErrorMessage(t *testing.T) {
	err := SecurityError{Message: "security failed"}

	if err.Error() != "security failed" {
		t.Fatalf("expected security failed, got %q", err.Error())
	}
}

func signToken(t *testing.T, secret string, claims AccessClaims) string {
	t.Helper()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return token
}

func validRegisteredClaims(subject string) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Subject:   subject,
		Issuer:    "issuer",
		Audience:  jwt.ClaimStrings{"audience"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
}
