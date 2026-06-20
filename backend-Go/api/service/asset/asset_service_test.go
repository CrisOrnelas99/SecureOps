package service

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/model"
	baserepository "secureops/backend-go/api/repository"
	baseservice "secureops/backend-go/api/service"
)

func TestAssetService(t *testing.T) {
	repo := &fakeAssetRepository{asset: sampleAsset(), assets: []model.Asset{sampleAsset()}}
	svc := NewAssetService(repo)
	ctx := newServiceContext(t, 42)

	if _, err := svc.GetAllAssets(ctx); err != nil {
		t.Fatalf("expected GetAllAssets to succeed, got %v", err)
	}
	if _, err := svc.CreateAsset(ctx, sampleAsset()); err != nil {
		t.Fatalf("expected CreateAsset to succeed, got %v", err)
	}
	if _, err := svc.UpdateAsset(ctx, 1, sampleAsset()); err != nil {
		t.Fatalf("expected UpdateAsset to succeed, got %v", err)
	}
}

func TestAssetServiceValidationAndTranslation(t *testing.T) {
	svc := NewAssetService(&fakeAssetRepository{findErr: baserepository.ErrAssetNotFound})
	ctx := newServiceContext(t, 42)

	if _, err := svc.GetAsset(ctx, 1); !errors.Is(err, baseservice.ErrNotFound) {
		t.Fatalf("expected not found translation, got %v", err)
	}
	if _, err := svc.CreateAsset(ctx, model.Asset{}); !errors.Is(err, baseservice.ErrInvalidRequestData) {
		t.Fatalf("expected invalid request data, got %v", err)
	}
}

type fakeAssetRepository struct {
	assets  []model.Asset
	asset   model.Asset
	findErr error
}

func (f *fakeAssetRepository) FindAllByUser(ec *appcontext.GinContext, userID int64) ([]model.Asset, error) {
	return f.assets, f.findErr
}
func (f *fakeAssetRepository) FindByIDForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Asset, error) {
	if f.findErr != nil {
		return model.Asset{}, f.findErr
	}
	return f.asset, nil
}
func (f *fakeAssetRepository) Save(ec *appcontext.GinContext, asset model.Asset) (model.Asset, error) {
	return asset, nil
}
func (f *fakeAssetRepository) UpdateForUser(ec *appcontext.GinContext, id int64, userID int64, asset model.Asset) (model.Asset, error) {
	return asset, nil
}
func (f *fakeAssetRepository) DeleteForUser(ec *appcontext.GinContext, id int64, userID int64) (model.Asset, error) {
	return f.asset, nil
}
func (f *fakeAssetRepository) AssignVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, error) {
	return f.asset, nil
}
func (f *fakeAssetRepository) RemoveVulnerabilityForUser(ec *appcontext.GinContext, assetID int64, userID int64, vulnerabilityID int64) (model.Asset, error) {
	return f.asset, nil
}

var _ baserepository.AssetRepository = (*fakeAssetRepository)(nil)

func newServiceContext(t *testing.T, userID int64) *appcontext.GinContext {
	t.Helper()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ec := appcontext.NewGinContext(ctx, "txn-123", log.New(io.Discard, "", 0))
	ctx.Set("userID", userID)
	appcontext.SetGinContext(ctx, ec)
	return ec
}

func sampleAsset() model.Asset {
	return model.Asset{Name: "Asset 1", Type: "Server", IPAddress: "10.0.0.10", Owner: "IT", Criticality: "High"}
}
