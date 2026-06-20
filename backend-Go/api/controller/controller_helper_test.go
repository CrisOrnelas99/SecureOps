package controller

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	appcontext "secureops/backend-go/api/context"
	"secureops/backend-go/api/dto"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestParseID(t *testing.T) {
	tests := []struct {
		input    string
		want     int64
		wantErr  bool
	}{
		{"1", 1, false},
		{"42", 42, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseID(tt.input)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestBindJSON(t *testing.T) {
	ec, recorder := newControllerContext(t, http.MethodPost, "/assets", `{"name":"Asset 1","type":"Server","ipAddress":"192.168.1.10","owner":"IT","criticality":"High"}`)
	ec.Request.Header.Set("Content-Type", "application/json")

	var request dto.AssetRequest
	if handled := BindJSON(ec, &request); handled {
		t.Fatal("expected request to bind")
	}
	if request.Name != "Asset 1" {
		t.Fatalf("expected name Asset 1, got %q", request.Name)
	}
	if recorder.Code != 200 {
		t.Fatalf("expected no error response, got %d", recorder.Code)
	}

	ec, _ = newControllerContext(t, http.MethodPost, "/assets", `{"name":"Asset 1","unknown":true}`)
	ec.Request.Header.Set("Content-Type", "application/json")
	if !BindJSON(ec, &request) {
		t.Fatal("expected unknown field to be rejected")
	}

	ec, _ = newControllerContext(t, http.MethodPost, "/assets", `{"name":"Asset 1"}{"name":"Asset 2"}`)
	ec.Request.Header.Set("Content-Type", "application/json")
	if !BindJSON(ec, &request) {
		t.Fatal("expected multiple JSON objects to be rejected")
	}

	ec, _ = newControllerContext(t, http.MethodPost, "/assets", `{"name":"Asset 1"}`)
	if !BindJSON(ec, &request) {
		t.Fatal("expected missing content type to be rejected")
	}
}

func TestHandleError(t *testing.T) {
	ec, recorder := newControllerContext(t, http.MethodGet, "/resource", "")
	if !HandleError(ec, http.StatusBadRequest, errors.New("boom"), "Invalid request body") {
		t.Fatal("expected error to be handled")
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response dto.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if response.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Code)
	}
	if response.Message != "Invalid request body" {
		t.Fatalf("expected message to match, got %q", response.Message)
	}
}

func newControllerContext(t *testing.T, method string, target string, body string) (*appcontext.GinContext, *httptest.ResponseRecorder) {
	t.Helper()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	ctx.Request = httptest.NewRequest(method, target, reader)
	if body == "" {
		ctx.Request.Body = http.NoBody
	}
	ec := appcontext.NewGinContext(ctx, "txn-123", log.New(io.Discard, "", 0))
	appcontext.SetGinContext(ctx, ec)
	return ec, recorder
}
