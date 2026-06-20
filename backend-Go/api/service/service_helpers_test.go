package service

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
	"secureops/backend-go/api/model"
	"secureops/backend-go/api/repository"
)

func TestServiceSentinels(t *testing.T) {
	if ErrInvalidRequestData.Error() != "invalid request data" {
		t.Fatal("unexpected sentinel message")
	}
}

func TestTranslateRepositoryError(t *testing.T) {
	if !errors.Is(TranslateRepositoryError(repository.ErrAssetNotFound), ErrNotFound) {
		t.Fatal("expected not found translation")
	}
	if !errors.Is(TranslateRepositoryError(repository.ErrInvalidData), ErrInvalidRequestData) {
		t.Fatal("expected invalid data translation")
	}
}

func TestValidateHelpers(t *testing.T) {
	if err := ValidateAsset(model.Asset{Name: "Asset 1", Type: "Server", IPAddress: "10.0.0.10", Owner: "IT", Criticality: "High"}); err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}
	if err := ValidateVulnerability(model.Vulnerability{CVEID: "CVE-1", Title: "Issue", Severity: "High", Description: "desc", Status: "Open"}); err != nil {
		t.Fatalf("expected valid vulnerability, got %v", err)
	}
	if err := ValidateRegisterRequest(dto.RegisterRequest{Username: "analyst", Email: "analyst@example.com", Password: "Password1!"}); err != nil {
		t.Fatalf("expected valid register request, got %v", err)
	}
}

func TestAuthenticatedUserID(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	ec := appcontext.NewGinContext(ctx, "txn", nil)
	ctx.Set("userID", int64(7))
	if id, err := AuthenticatedUserID(ec); err != nil || id != 7 {
		t.Fatalf("expected user id 7, got %d err=%v", id, err)
	}
	if _, err := AuthenticatedUserID(nil); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden for nil context, got %v", err)
	}
}

func TestNormalizeRegisterRequest(t *testing.T) {
	request := NormalizeRegisterRequest(dto.RegisterRequest{Username: " analyst ", Email: " ANALYST@EXAMPLE.COM "})
	if request.Username != "analyst" || request.Email != "analyst@example.com" {
		t.Fatalf("unexpected normalized request: %#v", request)
	}
}

var _ = gorm.ErrRecordNotFound
