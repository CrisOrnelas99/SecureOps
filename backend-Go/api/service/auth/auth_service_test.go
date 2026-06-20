package service

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
	"secureops/backend-go/api/model"
	baserepository "secureops/backend-go/api/repository"
	baseservice "secureops/backend-go/api/service"
	"secureops/backend-go/api/security"
)

func TestAuthService(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("Password1!"), bcrypt.DefaultCost)
	repo := &fakeUserRepository{
		user: model.User{ID: 1, Username: "analyst", Email: "analyst@example.com", PasswordHash: string(hash), Role: model.RoleUser},
	}
	svc := NewAuthService(security.NewJWTManager("test-secret", time.Hour, "issuer", "audience"), repo)
	ctx := newAuthServiceContext(t)

	if err := svc.Register(ctx, dto.RegisterRequest{Username: "analyst", Email: "analyst@example.com", Password: "Password1!"}); err != nil {
		t.Fatalf("expected Register to succeed, got %v", err)
	}
	response, err := svc.Login(ctx, dto.LoginRequest{UserOrEmail: "analyst", Password: "Password1!"})
	if err != nil {
		t.Fatalf("expected Login to succeed, got %v", err)
	}
	if response.Token == "" {
		t.Fatal("expected token to be populated")
	}
}

func TestAuthServiceValidationAndTranslation(t *testing.T) {
	ctx := newAuthServiceContext(t)
	svc := NewAuthService(security.NewJWTManager("test-secret", time.Hour, "issuer", "audience"), &fakeUserRepository{findErr: gorm.ErrRecordNotFound})

	if err := svc.Register(ctx, dto.RegisterRequest{Username: "ab", Email: "bad", Password: "short"}); !errors.Is(err, baseservice.ErrInvalidRequestData) {
		t.Fatalf("expected invalid request data, got %v", err)
	}
	if _, err := svc.Login(ctx, dto.LoginRequest{UserOrEmail: "missing", Password: "Password1!"}); !errors.Is(err, baseservice.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

type fakeUserRepository struct {
	user     model.User
	findErr  error
	exists   bool
}

func (f *fakeUserRepository) ExistsByUsername(ec *appcontext.GinContext, username string) (bool, error) { return f.exists, nil }
func (f *fakeUserRepository) ExistsByEmail(ec *appcontext.GinContext, email string) (bool, error)    { return f.exists, nil }
func (f *fakeUserRepository) Save(ec *appcontext.GinContext, user model.User) error                   { return nil }
func (f *fakeUserRepository) FindByUsernameOrEmail(ec *appcontext.GinContext, userOrEmail string) (model.User, error) {
	return f.user, f.findErr
}
func (f *fakeUserRepository) FindByUsername(ec *appcontext.GinContext, username string) (model.User, error) {
	return f.user, f.findErr
}

var _ baserepository.UserRepository = (*fakeUserRepository)(nil)

func newAuthServiceContext(t *testing.T) *appcontext.GinContext {
	t.Helper()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ec := appcontext.NewGinContext(ctx, "txn-123", log.New(io.Discard, "", 0))
	appcontext.SetGinContext(ctx, ec)
	return ec
}
