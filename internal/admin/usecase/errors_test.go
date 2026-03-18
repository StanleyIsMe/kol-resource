package usecase

import (
	"errors"
	"net/http"
	"testing"
)

func TestInternalServerError(t *testing.T) {
	t.Parallel()

	err := InternalServerError{err: errors.New("something went wrong")}

	if code := err.ErrorCode(); code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("ErrorCode() = %q, want %q", code, "INTERNAL_SERVER_ERROR")
	}

	if msg := err.ErrorMsg(); msg != "Internal Server Error" {
		t.Errorf("ErrorMsg() = %q, want %q", msg, "Internal Server Error")
	}

	if errStr := err.Error(); errStr != "internal server error: something went wrong" {
		t.Errorf("Error() = %q, want %q", errStr, "internal server error: something went wrong")
	}

	if status := err.HTTPStatusCode(); status != http.StatusInternalServerError {
		t.Errorf("HTTPStatusCode() = %d, want %d", status, http.StatusInternalServerError)
	}
}

func TestDumplicatedUsernameError(t *testing.T) {
	t.Parallel()

	err := DumplicatedUsernameError{username: "admin"}

	if code := err.ErrorCode(); code != "DUPLICATED_USERNAME" {
		t.Errorf("ErrorCode() = %q, want %q", code, "DUPLICATED_USERNAME")
	}

	if msg := err.ErrorMsg(); msg != "username admin already exists" {
		t.Errorf("ErrorMsg() = %q, want %q", msg, "username admin already exists")
	}

	if errStr := err.Error(); errStr != "username admin already exists" {
		t.Errorf("Error() = %q, want %q", errStr, "username admin already exists")
	}

	if status := err.HTTPStatusCode(); status != http.StatusBadRequest {
		t.Errorf("HTTPStatusCode() = %d, want %d", status, http.StatusBadRequest)
	}
}

func TestUnauthorizedError(t *testing.T) {
	t.Parallel()

	err := UnauthorizedError{err: errors.New("invalid token")}

	if code := err.ErrorCode(); code != "UNAUTHORIZED" {
		t.Errorf("ErrorCode() = %q, want %q", code, "UNAUTHORIZED")
	}

	if msg := err.ErrorMsg(); msg != "Unauthorized" {
		t.Errorf("ErrorMsg() = %q, want %q", msg, "Unauthorized")
	}

	if errStr := err.Error(); errStr != "unauthorized: invalid token" {
		t.Errorf("Error() = %q, want %q", errStr, "unauthorized: invalid token")
	}

	if status := err.HTTPStatusCode(); status != http.StatusUnauthorized {
		t.Errorf("HTTPStatusCode() = %d, want %d", status, http.StatusUnauthorized)
	}
}
