package email

import (
	"flag"
	"os"
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

func TestEmailJobStatus_CanCancel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status EmailJobStatus
		want   bool
	}{
		{name: "pending can cancel", status: EmailJobStatusPending, want: true},
		{name: "processing can cancel", status: EmailJobStatusProcessing, want: true},
		{name: "success cannot cancel", status: EmailJobStatusSuccess, want: false},
		{name: "canceled cannot cancel", status: EmailJobStatusCanceled, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.status.CanCancel(); got != tt.want {
				t.Errorf("EmailJobStatus(%q).CanCancel() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestEmailJobStatus_CanStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status EmailJobStatus
		want   bool
	}{
		{name: "canceled can start", status: EmailJobStatusCanceled, want: true},
		{name: "pending cannot start", status: EmailJobStatusPending, want: false},
		{name: "success cannot start", status: EmailJobStatusSuccess, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.status.CanStart(); got != tt.want {
				t.Errorf("EmailJobStatus(%q).CanStart() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestEmailJobStatus_ToPointer(t *testing.T) {
	t.Parallel()

	status := EmailJobStatusPending
	ptr := status.ToPointer()

	if ptr == nil {
		t.Fatal("ToPointer() returned nil")
	}

	if *ptr != status {
		t.Errorf("*ToPointer() = %q, want %q", *ptr, status)
	}
}

func TestEmailLogStatus_ToPointer(t *testing.T) {
	t.Parallel()

	status := EmailLogStatusPending
	ptr := status.ToPointer()

	if ptr == nil {
		t.Fatal("ToPointer() returned nil")
	}

	if *ptr != status {
		t.Errorf("*ToPointer() = %q, want %q", *ptr, status)
	}
}
