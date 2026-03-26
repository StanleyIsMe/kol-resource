package kol

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

func TestSex_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sex  Sex
		want bool
	}{
		{name: "male is valid", sex: "m", want: true},
		{name: "female is valid", sex: "f", want: true},
		{name: "unknown is invalid", sex: "x", want: false},
		{name: "empty is invalid", sex: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.sex.IsValid(); got != tt.want {
				t.Errorf("Sex(%q).IsValid() = %v, want %v", tt.sex, got, tt.want)
			}
		})
	}
}

func TestSex_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sex  Sex
		want string
	}{
		{name: "male string", sex: SexMale, want: "m"},
		{name: "female string", sex: SexFemale, want: "f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.sex.String(); got != tt.want {
				t.Errorf("Sex(%q).String() = %q, want %q", tt.sex, got, tt.want)
			}
		})
	}
}
