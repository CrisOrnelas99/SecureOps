package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestLoadUsesDefaults(t *testing.T) {
	clearConfigEnv(t)

	cfg := Load()

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://secureops:secureops@localhost:5432/secureops" {
		t.Fatalf("unexpected default database URL: %q", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "" {
		t.Fatalf("expected default JWT secret to be empty, got %q", cfg.JWTSecret)
	}
	if cfg.JWTExpiration != time.Hour {
		t.Fatalf("expected default JWT expiration %s, got %s", time.Hour, cfg.JWTExpiration)
	}
	if cfg.JWTIssuer != "secureops-lite" {
		t.Fatalf("expected default JWT issuer secureops-lite, got %q", cfg.JWTIssuer)
	}
	if cfg.JWTAudience != "secureops-lite-api" {
		t.Fatalf("expected default JWT audience secureops-lite-api, got %q", cfg.JWTAudience)
	}
}

func TestLoadUsesEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DB_HOST", "db")
	t.Setenv("POSTGRES_PORT", "15432")
	t.Setenv("POSTGRES_DB", "app")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "pass")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_ISSUER", "issuer")
	t.Setenv("JWT_AUDIENCE", "audience")
	t.Setenv("JWT_EXPIRATION_MS", "60000")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Fatalf("expected configured port 9090, got %q", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@db:15432/app" {
		t.Fatalf("unexpected configured database URL: %q", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "test-secret" {
		t.Fatalf("expected configured JWT secret, got %q", cfg.JWTSecret)
	}
	if cfg.JWTExpiration != time.Minute {
		t.Fatalf("expected configured JWT expiration %s, got %s", time.Minute, cfg.JWTExpiration)
	}
	if cfg.JWTIssuer != "issuer" {
		t.Fatalf("expected configured JWT issuer issuer, got %q", cfg.JWTIssuer)
	}
	if cfg.JWTAudience != "audience" {
		t.Fatalf("expected configured JWT audience audience, got %q", cfg.JWTAudience)
	}
}

func TestLoadFallsBackForInvalidJWTExpiration(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("JWT_EXPIRATION_MS", "not-a-number")

	cfg := Load()

	if cfg.JWTExpiration != time.Hour {
		t.Fatalf("expected invalid JWT expiration to fall back to %s, got %s", time.Hour, cfg.JWTExpiration)
	}
}

func TestLoadFallsBackForNonPositiveJWTExpiration(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("JWT_EXPIRATION_MS", "0")

	cfg := Load()

	if cfg.JWTExpiration != time.Hour {
		t.Fatalf("expected non-positive JWT expiration to fall back to %s, got %s", time.Hour, cfg.JWTExpiration)
	}
}

func TestPasswordCostUsesBcryptDefault(t *testing.T) {
	if PasswordCost() != bcrypt.DefaultCost {
		t.Fatalf("expected password cost %d, got %d", bcrypt.DefaultCost, PasswordCost())
	}
}

func TestCorsConfigSetsHeadersAndContinuesForNormalRequests(t *testing.T) {
	router := gin.New()
	router.Use(CorsConfig())
	router.GET("/resource", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/resource", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	assertCorsHeaders(t, recorder)
}

func TestCorsConfigAbortsOptionsPreflight(t *testing.T) {
	router := gin.New()
	router.Use(CorsConfig())

	handlerCalled := false
	router.OPTIONS("/resource", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/resource", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
	if handlerCalled {
		t.Fatal("expected OPTIONS request to abort before route handler")
	}
	assertCorsHeaders(t, recorder)
}

func assertCorsHeaders(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "http://localhost:4200",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE",
		"Access-Control-Allow-Headers":     "Authorization, Content-Type",
		"Access-Control-Allow-Credentials": "true",
	}

	for header, expected := range expectedHeaders {
		if actual := recorder.Header().Get(header); actual != expected {
			t.Fatalf("expected %s header %q, got %q", header, expected, actual)
		}
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()

	keys := []string{
		"PORT",
		"DB_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"JWT_SECRET",
		"JWT_ISSUER",
		"JWT_AUDIENCE",
		"JWT_EXPIRATION_MS",
	}

	for _, key := range keys {
		t.Setenv(key, "")
	}
}
