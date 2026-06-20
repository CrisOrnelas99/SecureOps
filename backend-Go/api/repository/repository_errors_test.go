package repository

import "testing"

func TestRepositoryErrorMessages(t *testing.T) {
	if ErrReadFailed.Error() != "read failed" {
		t.Fatal("unexpected read failed message")
	}
	if ErrDuplicateAssignment.Error() != "duplicate asset vulnerability assignment" {
		t.Fatal("unexpected duplicate assignment message")
	}
}
