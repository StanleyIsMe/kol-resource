package business

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestInternalServerError(t *testing.T) {
	t.Parallel()

	inner := errors.New("something broke")
	err := InternalServerError{err: inner}

	if code := err.ErrorCode(); code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("ErrorCode() = %q, want %q", code, "INTERNAL_SERVER_ERROR")
	}

	if msg := err.ErrorMsg(); msg != "Internal Server Error" {
		t.Errorf("ErrorMsg() = %q, want %q", msg, "Internal Server Error")
	}

	if errStr := err.Error(); !strings.Contains(errStr, "something broke") {
		t.Errorf("Error() = %q, expected to contain inner error message", errStr)
	}

	if status := err.HTTPStatusCode(); status != http.StatusInternalServerError {
		t.Errorf("HTTPStatusCode() = %d, want %d", status, http.StatusInternalServerError)
	}
}
