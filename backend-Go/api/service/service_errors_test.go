// Package service verifies service-level sentinel errors.
package service

import "testing"

// TestServiceSentinels verifies the core service sentinel messages.
func TestServiceSentinels(t *testing.T) {
	if ErrInvalidRequestData.Error() != "invalid request data" {
		t.Fatal("unexpected sentinel message")
	}
}
