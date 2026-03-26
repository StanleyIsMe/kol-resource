package usecase

import (
	"net/http"
	"strings"
	"testing"
)

func TestDuplicatedResourceError(t *testing.T) {
	t.Parallel()

	err := DuplicatedResourceError{resource: "Kol", name: "test-kol"}

	if code := err.ErrorCode(); code != "DUPLICATED_RESOURCE" {
		t.Errorf("ErrorCode() = %q, want %q", code, "DUPLICATED_RESOURCE")
	}

	if msg := err.ErrorMsg(); !strings.Contains(msg, "Kol") || !strings.Contains(msg, "test-kol") {
		t.Errorf("ErrorMsg() = %q, expected to contain resource and name", msg)
	}

	if errStr := err.Error(); !strings.Contains(errStr, "Kol") || !strings.Contains(errStr, "test-kol") {
		t.Errorf("Error() = %q, expected to contain resource and name", errStr)
	}

	if status := err.HTTPStatusCode(); status != http.StatusBadRequest {
		t.Errorf("HTTPStatusCode() = %d, want %d", status, http.StatusBadRequest)
	}
}

func TestNotFoundError(t *testing.T) {
	t.Parallel()

	err := NotFoundError{resource: "Kol", id: "some-id"}

	if code := err.ErrorCode(); code != "NOT_FOUND" {
		t.Errorf("ErrorCode() = %q, want %q", code, "NOT_FOUND")
	}

	if msg := err.ErrorMsg(); !strings.Contains(msg, "Kol") || !strings.Contains(msg, "some-id") {
		t.Errorf("ErrorMsg() = %q, expected to contain resource and id", msg)
	}

	if errStr := err.Error(); !strings.Contains(errStr, "Kol") || !strings.Contains(errStr, "some-id") {
		t.Errorf("Error() = %q, expected to contain resource and id", errStr)
	}

	if status := err.HTTPStatusCode(); status != http.StatusNotFound {
		t.Errorf("HTTPStatusCode() = %d, want %d", status, http.StatusNotFound)
	}
}
