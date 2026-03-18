package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestErrDataNotFound(t *testing.T) {
	t.Parallel()

	if ErrDataNotFound.Error() != "data not found" {
		t.Errorf("ErrDataNotFound = %q, want %q", ErrDataNotFound.Error(), "data not found")
	}
}

func TestGenerateUUIDError_Error(t *testing.T) {
	t.Parallel()

	err := GenerateUUIDError{Err: errors.New("rand failure")}

	got := err.Error()
	if !strings.Contains(got, "rand failure") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestInsertRecordError_Error(t *testing.T) {
	t.Parallel()

	err := InsertRecordError{Err: errors.New("insert failed")}

	got := err.Error()
	if !strings.Contains(got, "insert failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestQueryRecordError_Error(t *testing.T) {
	t.Parallel()

	err := QueryRecordError{Err: errors.New("query failed")}

	got := err.Error()
	if !strings.Contains(got, "query failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestUpdateRecordError_Error(t *testing.T) {
	t.Parallel()

	err := UpdateRecordError{Err: errors.New("update failed")}

	got := err.Error()
	if !strings.Contains(got, "update failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestDeleteRecordError_Error(t *testing.T) {
	t.Parallel()

	err := DeleteRecordError{Err: errors.New("delete failed")}

	got := err.Error()
	if !strings.Contains(got, "delete failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestUUIDInvalidError_Error(t *testing.T) {
	t.Parallel()

	err := UUIDInvalidError{Field: "kol_id", UUID: "not-a-uuid"}

	got := err.Error()
	if !strings.Contains(got, "kol_id") || !strings.Contains(got, "not-a-uuid") {
		t.Errorf("Error() = %q, expected to contain field and uuid", got)
	}
}

func TestSexInvalidError_Error(t *testing.T) {
	t.Parallel()

	err := SexInvalidError{Sex: "x"}

	got := err.Error()
	if !strings.Contains(got, "x") {
		t.Errorf("Error() = %q, expected to contain sex value", got)
	}
}
