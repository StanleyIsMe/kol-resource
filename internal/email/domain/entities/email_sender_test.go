package entities

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

func TestEmailSender_SetDailyRateLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rateLimit int
		want      int
	}{
		{name: "positive value sets rate limit", rateLimit: 100, want: 100},
		{name: "zero sets default", rateLimit: 0, want: DefaultDaRateLimit},
		{name: "negative sets default", rateLimit: -1, want: DefaultDaRateLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sender := &EmailSender{}
			sender.SetDailyRateLimit(tt.rateLimit)

			if sender.RateLimit != tt.want {
				t.Errorf("RateLimit = %d, want %d", sender.RateLimit, tt.want)
			}
		})
	}
}
