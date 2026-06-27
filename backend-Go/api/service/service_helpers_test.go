// Package service verifies shared service helper behavior.
package service

import (
	"errors"
	"testing"
)

// TestCVEIDValidation verifies strict CVE ID allowlist behavior.
func TestCVEIDValidation(t *testing.T) {
	if NormalizeCVEID(" cve-2021-44228 ") != "CVE-2021-44228" {
		t.Fatal("expected CVE ID to be trimmed and uppercased")
	}

	valid := []string{
		"CVE-2021-44228",
		"cve-2024-12345",
	}
	for _, cveID := range valid {
		if err := ValidateCVEID(cveID); err != nil {
			t.Fatalf("expected %q to be valid, got %v", cveID, err)
		}
	}

	invalid := []string{
		"CVE-21-44228",
		"CVE-2021-123",
		"https://nvd.nist.gov/vuln/detail/CVE-2021-44228",
		"CVE-2021-44228?redirect=https://example.com",
	}
	for _, cveID := range invalid {
		if !errors.Is(ValidateCVEID(cveID), ErrInvalidRequestData) {
			t.Fatalf("expected %q to be rejected", cveID)
		}
	}
}
