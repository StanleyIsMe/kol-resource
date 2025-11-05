package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"time"

	"github.com/rs/zerolog"
)

const (
	defaultEmailScheduleInterval = 1 * time.Minute
	defaultEmailCountPerMinute   = 4
)

type EmailSchedule struct {
	repo      domain.Repository
	emailRepo domain.EmailRepository
	interval  time.Duration
}

func NewEmailSchedule(repo domain.Repository, emailRepo domain.EmailRepository, interval time.Duration) *EmailSchedule {
	return &EmailSchedule{repo: repo, emailRepo: emailRepo, interval: interval}
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

	// emailSenders, err := s.repo.AllEmailSenders(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to get email senders: %w", err)
	// }

	// calculate 24h period	 emails count
	// get pending email job

	emailJobs, err := s.repo.GrabEmailJob(ctx)
	if err != nil {
		return fmt.Errorf("repo.GrabEmailJob error: %w", err)
	}

	for _, emailJob := range emailJobs {
		// get email sender by sender_id
		emailSender, err := s.repo.GetEmailSenderByID(ctx, emailJob.SenderID)
		if err != nil {
			if errors.Is(err, commonErrors.ErrDataNotFound) {
				return commonErrors.NotFoundError{
					Resource: "email sender",
					ID:       emailJob.SenderID,
				}
			}

			return fmt.Errorf("repo.GetEmailSenderByID error: %w", err)
		}

		sentEmailsCount, err := s.repo.CountSentEmailsLast24Hours(ctx, emailSender.ID)
		if err != nil {
			return fmt.Errorf("repo.CountSentEmailsLast24Hours error: %w", err)
		}

		if sentEmailsCount >= int64(float64(emailSender.RateLimit)*0.9) {
			zerolog.Ctx(ctx).Info().Fields(map[string]any{
				"sender_id":         emailSender.ID,
				"sent_emails_count": sentEmailsCount,
				"rate_limit":        emailSender.RateLimit,
			}).Msg("email sender rate limit exceeded")

			continue
		}

		// get email receiver from email log
		emailLog, err := s.repo.GrabPendingEmailLogByJobID(ctx, emailJob.ID)
		if err != nil {
			return fmt.Errorf("repo.GrabPendingEmailLogByJobID error: %w", err)
		}

		memo := ""
		// send email
		var emailPayload domain.SendEmailParams
		err = json.Unmarshal([]byte(emailJob.Payload), &emailPayload)
		if err != nil {
			return fmt.Errorf("json.Unmarshal error: %w", err)
		}
		err = s.emailRepo.SendEmail(ctx, domain.SendEmailParams{
			Subject: emailPayload.Subject,
			Body:    emailPayload.Body,
			ToEmails: []domain.ToEmail{
				{
					Email: emailLog.Email,
					Name:  emailLog.KolName,
				},
			},
			Images: emailPayload.Images,
		})
		if err != nil {
			zerolog.Ctx(ctx).Error().Fields(map[string]any{
				"error":        err,
				"email_log_id": emailLog.ID,
			}).Msg("emailRepo.SendEmail error")

			memo = fmt.Sprintf("failed to send email: %s", err.Error())
			
			return fmt.Errorf("emailRepo.SendEmail error: %w", err)
		}

		// update email log information
		err = s.repo.UpdateEmailLog(ctx, domain.UpdateEmailLogParam{
			ID:     emailLog.ID,
			Status: email.EmailLogStatusFailed.ToPointer(),
			Memo:   memo,
		})
		if err != nil {
			return fmt.Errorf("repo.UpdateEmailLog error: %w", err)
		}

		// update email job information
	}

	return nil
}
