package main

import (
	"context"
	"strings"
	"testing"

	"secureops/backend-go/api/config"
)

func TestRunBootstrapSkipsWhenDisabled(t *testing.T) {
	cfg := config.Config{
		Environment:      "development",
		BootstrapDevData: false,
	}

	if err := runBootstrap(context.Background(), nil, cfg); err != nil {
		t.Fatalf("expected disabled bootstrap to skip without error, got %v", err)
	}
}

func TestRunBootstrapRejectsProduction(t *testing.T) {
	cfg := config.Config{
		Environment:      "production",
		BootstrapDevData: true,
	}

	err := runBootstrap(context.Background(), nil, cfg)
	if err == nil {
		t.Fatal("expected production bootstrap to fail")
	}
	if !strings.Contains(err.Error(), "production") {
		t.Fatalf("expected production error, got %v", err)
	}
}
