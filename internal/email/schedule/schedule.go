package schedule

import (
	"context"
	"fmt"
	"kolresource/internal/email/domain"
	"time"

	"github.com/rs/zerolog"
)

const (
	defaultEmailScheduleInterval = 1 * time.Minute
	defaultEmailCountPerMinute   = 4
)

type EmailSchedule struct {
	repo     domain.Repository
	interval time.Duration
}

func NewEmailSchedule(repo domain.Repository, interval time.Duration) *EmailSchedule {
	return &EmailSchedule{repo: repo, interval: interval}
}

func (s *EmailSchedule) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.sendEmail(ctx); err != nil {
					zerolog.Ctx(ctx).Error().Fields(map[string]any{
						"error": err,
					}).Msg("email schedule send email error")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *EmailSchedule) sendEmail(ctx context.Context) error {
	// list email senders

	emailSenders, err := s.repo.AllEmailSenders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get email senders: %w", err)
	}

	// calculate 24h period	 emails count
	// get pending email job
	return nil
}
