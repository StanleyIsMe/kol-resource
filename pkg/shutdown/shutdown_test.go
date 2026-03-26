package shutdown

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/rs/zerolog"
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

func TestNew(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()

	s := New(&logger)

	if s == nil {
		t.Fatal("New() returned nil")
	}

	if len(s.Hooks()) != 0 {
		t.Errorf("New() hooks count = %d, want 0", len(s.Hooks()))
	}

	if s.gracePeriodDuration != defaultGracePeriodDuration {
		t.Errorf("New() gracePeriodDuration = %v, want %v", s.gracePeriodDuration, defaultGracePeriodDuration)
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	s := New(&logger)

	s.Add("hook1", func(ctx context.Context) error { return nil })
	s.Add("hook2", func(ctx context.Context) error { return nil })

	hooks := s.Hooks()
	if len(hooks) != 2 {
		t.Fatalf("Hooks() count = %d, want 2", len(hooks))
	}

	if hooks[0].Name != "hook1" {
		t.Errorf("Hooks()[0].Name = %q, want %q", hooks[0].Name, "hook1")
	}

	if hooks[1].Name != "hook2" {
		t.Errorf("Hooks()[1].Name = %q, want %q", hooks[1].Name, "hook2")
	}
}

func TestWithHooks(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()

	hooks := []Hook{
		{Name: "db", ShutdownFn: func(ctx context.Context) error { return nil }},
		{Name: "cache", ShutdownFn: func(ctx context.Context) error { return nil }},
	}

	s := New(&logger, WithHooks(hooks))

	got := s.Hooks()
	if len(got) != 2 {
		t.Fatalf("Hooks() count = %d, want 2", len(got))
	}

	if got[0].Name != "db" {
		t.Errorf("Hooks()[0].Name = %q, want %q", got[0].Name, "db")
	}

	if got[1].Name != "cache" {
		t.Errorf("Hooks()[1].Name = %q, want %q", got[1].Name, "cache")
	}
}

func TestWithGracePeriodDuration(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	duration := 10 * time.Second

	s := New(&logger, WithGracePeriodDuration(duration))

	if s.gracePeriodDuration != duration {
		t.Errorf("gracePeriodDuration = %v, want %v", s.gracePeriodDuration, duration)
	}
}

func TestListen(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(io.Discard)
	hookCalled := false
	s := New(&logger, WithGracePeriodDuration(5*time.Second))
	s.Add("test", func(ctx context.Context) error {
		hookCalled = true
		return nil
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()

	err := s.Listen(context.Background(), syscall.SIGINT)
	if err != nil {
		t.Errorf("Listen() error = %v", err)
	}
	if !hookCalled {
		t.Error("hook was not called")
	}
}

func TestListen_WithError(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(io.Discard)
	s := New(&logger, WithGracePeriodDuration(5*time.Second))
	s.Add("failing-hook", func(ctx context.Context) error {
		return errors.New("shutdown failed")
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()

	err := s.Listen(context.Background(), syscall.SIGINT)
	if err == nil {
		t.Error("Listen() expected error, got nil")
	}
}
