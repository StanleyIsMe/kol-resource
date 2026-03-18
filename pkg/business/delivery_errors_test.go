package business

import (
	"errors"
	"net/http"
	"testing"
)

func TestUseCaesErrorToErrorResp_WithUseCaseError(t *testing.T) {
	t.Parallel()

	inner := errors.New("something broke")
	ucErr := InternalServerError{err: inner}

	statusCode, resp := UseCaesErrorToErrorResp(ucErr)

	if statusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", statusCode, http.StatusInternalServerError)
	}

	if resp.ErrorCode != "INTERNAL_SERVER_ERROR" {
		t.Errorf("ErrorCode = %q, want %q", resp.ErrorCode, "INTERNAL_SERVER_ERROR")
	}

	if resp.ErrorMessage != "Internal Server Error" {
		t.Errorf("ErrorMessage = %q, want %q", resp.ErrorMessage, "Internal Server Error")
	}
}

func TestUseCaesErrorToErrorResp_WithNonUseCaseError(t *testing.T) {
	t.Parallel()

	plainErr := errors.New("plain error")

	statusCode, resp := UseCaesErrorToErrorResp(plainErr)

	if statusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", statusCode, http.StatusInternalServerError)
	}

	if resp.ErrorCode != "INTERNAL_SERVER_ERROR" {
		t.Errorf("ErrorCode = %q, want %q", resp.ErrorCode, "INTERNAL_SERVER_ERROR")
	}

	if resp.ErrorMessage != "Internal Server Error" {
		t.Errorf("ErrorMessage = %q, want %q", resp.ErrorMessage, "Internal Server Error")
	}
}
