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

func TestJobStatus_CanCancel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status JobStatus
		want   bool
	}{
		{name: "pending can cancel", status: JobStatusPending, want: true},
		{name: "processing can cancel", status: JobStatusProcessing, want: true},
		{name: "success cannot cancel", status: JobStatusSuccess, want: false},
		{name: "canceled cannot cancel", status: JobStatusCanceled, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.status.CanCancel(); got != tt.want {
				t.Errorf("JobStatus(%q).CanCancel() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestJobStatus_CanStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status JobStatus
		want   bool
	}{
		{name: "canceled can start", status: JobStatusCanceled, want: true},
		{name: "pending cannot start", status: JobStatusPending, want: false},
		{name: "success cannot start", status: JobStatusSuccess, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.status.CanStart(); got != tt.want {
				t.Errorf("JobStatus(%q).CanStart() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestJobStatus_ToPointer(t *testing.T) {
	t.Parallel()

	status := JobStatusPending
	ptr := status.ToPointer()

	if ptr == nil {
		t.Fatal("ToPointer() returned nil")
	}

	if *ptr != status {
		t.Errorf("*ToPointer() = %q, want %q", *ptr, status)
	}
}

func TestLogStatus_ToPointer(t *testing.T) {
	t.Parallel()

	status := LogStatusPending
	ptr := status.ToPointer()

	if ptr == nil {
		t.Fatal("ToPointer() returned nil")
	}

	if *ptr != status {
		t.Errorf("*ToPointer() = %q, want %q", *ptr, status)
	}
}
