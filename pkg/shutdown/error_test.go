package shutdown

import (
	"errors"
	"strings"
	"testing"
)

func TestShutdownError_Empty(t *testing.T) {
	t.Parallel()

	sErr := newShutdownError()

	got := sErr.Error()
	if got != "" {
		t.Errorf("Error() = %q, want empty string", got)
	}
}

func TestShutdownError_WithErrors(t *testing.T) {
	t.Parallel()

	sErr := newShutdownError()
	sErr["db"] = errors.New("connection refused")
	sErr["cache"] = errors.New("timeout")

	got := sErr.Error()

	if !strings.Contains(got, "db err: connection refused") {
		t.Errorf("Error() = %q, want it to contain %q", got, "db err: connection refused")
	}

	if !strings.Contains(got, "cache err: timeout") {
		t.Errorf("Error() = %q, want it to contain %q", got, "cache err: timeout")
	}
}
