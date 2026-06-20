package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
	"secureops/backend-go/api/middleware"
	"secureops/backend-go/api/model"
	"secureops/backend-go/api/security"
)

func TestHealth(t *testing.T) {
	router := gin.New()
	router.GET("/api/health", Health)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()
	jwtManager := security.NewJWTManager("test-secret", time.Hour, "issuer", "audience")
	lookup := &fakeUserLookup{exists: true, user: model.User{ID: 1, Username: "analyst", Role: model.RoleUser}}

	authController := &fakeAuthController{}
	assetController := &fakeAssetController{}
	vulnerabilityController := &fakeVulnerabilityController{}

	RegisterRoutes(router, jwtManager, lookup, authController, assetController, vulnerabilityController)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected health endpoint to be registered, got %d", recorder.Code)
	}
}

type fakeUserLookup struct {
	exists bool
	user   model.User
}

func (f *fakeUserLookup) ExistsByUsername(ec *context.GinContext, username string) (bool, error) {
	return f.exists, nil
}

func (f *fakeUserLookup) FindByUsername(ec *context.GinContext, username string) (model.User, error) {
	return f.user, nil
}

type fakeAuthController struct{}

func (f *fakeAuthController) Register(ec *context.GinContext) {}
func (f *fakeAuthController) Login(ec *context.GinContext)    {}

type fakeAssetController struct{}

func (f *fakeAssetController) GetAssets(ec *context.GinContext)           {}
func (f *fakeAssetController) GetAsset(ec *context.GinContext)            {}
func (f *fakeAssetController) CreateAsset(ec *context.GinContext)         {}
func (f *fakeAssetController) UpdateAsset(ec *context.GinContext)         {}
func (f *fakeAssetController) DeleteAsset(ec *context.GinContext)         {}
func (f *fakeAssetController) AssignVulnerability(ec *context.GinContext) {}
func (f *fakeAssetController) RemoveVulnerability(ec *context.GinContext) {}

type fakeVulnerabilityController struct{}

func (f *fakeVulnerabilityController) GetVulnerabilities(ec *context.GinContext)  {}
func (f *fakeVulnerabilityController) GetVulnerability(ec *context.GinContext)    {}
func (f *fakeVulnerabilityController) CreateVulnerability(ec *context.GinContext) {}
func (f *fakeVulnerabilityController) UpdateVulnerability(ec *context.GinContext) {}
func (f *fakeVulnerabilityController) DeleteVulnerability(ec *context.GinContext) {}

var _ middleware.UserLookup = (*fakeUserLookup)(nil)

var _ = dto.ErrorResponse{}
