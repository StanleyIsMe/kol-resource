package errors

import (
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestGenerateUUIDError_Error(t *testing.T) {
	t.Parallel()

	err := GenerateUUIDError{Err: errors.New("rand failure")}

	if got := err.Error(); !strings.Contains(got, "rand failure") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestInsertRecordError_Error(t *testing.T) {
	t.Parallel()

	err := InsertRecordError{Err: errors.New("insert failed")}

	if got := err.Error(); !strings.Contains(got, "insert failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestQueryRecordError_Error(t *testing.T) {
	t.Parallel()

	err := QueryRecordError{Err: errors.New("query failed")}

	if got := err.Error(); !strings.Contains(got, "query failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestUpdateRecordError_Error(t *testing.T) {
	t.Parallel()

	err := UpdateRecordError{Err: errors.New("update failed")}

	if got := err.Error(); !strings.Contains(got, "update failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestDeleteRecordError_Error(t *testing.T) {
	t.Parallel()

	err := DeleteRecordError{Err: errors.New("delete failed")}

	if got := err.Error(); !strings.Contains(got, "delete failed") {
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

	if got := err.Error(); !strings.Contains(got, "x") {
		t.Errorf("Error() = %q, expected to contain sex value", got)
	}
}

func TestTransactionError_Error(t *testing.T) {
	t.Parallel()

	err := TransactionError{Err: errors.New("tx failed")}

	if got := err.Error(); !strings.Contains(got, "tx failed") {
		t.Errorf("Error() = %q, expected to contain inner error", got)
	}
}

func TestTransactionRollbackError_Error(t *testing.T) {
	t.Parallel()

	err := &TransactionRollbackError{
		DbErr: errors.New("db error"),
		TxErr: errors.New("rollback error"),
	}

	got := err.Error()
	if !strings.Contains(got, "db error") || !strings.Contains(got, "rollback error") {
		t.Errorf("Error() = %q, expected to contain both errors", got)
	}
}

func TestIsUniqueViolationError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "non_pgconn_error_returns_false",
			err:  errors.New("some random error"),
			want: false,
		},
		{
			name: "pgconn_error_wrong_code_returns_false",
			err:  &pgconn.PgError{Code: "42601"},
			want: false,
		},
		{
			name: "pgconn_error_unique_violation_returns_true",
			err:  &pgconn.PgError{Code: "23505"},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsUniqueViolationError(tt.err)
			if got != tt.want {
				t.Errorf("IsUniqueViolationError() = %v, want %v", got, tt.want)
			}
		})
	}
}
