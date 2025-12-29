package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	defaultEmailScheduleInterval = 1 * time.Minute
	defaultEmailCountPerMinute   = 2
)

type EmailSchedule struct {
	repo      domain.Repository
	emailRepo domain.EmailRepository
	interval  time.Duration
}

func NewEmailSchedule(repo domain.Repository, emailRepo domain.EmailRepository, interval time.Duration) *EmailSchedule {
	if interval <= 0 {
		interval = defaultEmailScheduleInterval
	}

	return &EmailSchedule{repo: repo, emailRepo: emailRepo, interval: interval}
}

func (s *EmailSchedule) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				jobCtx, jobCancel := context.WithTimeout(ctx, s.interval)
				jobCtx = zerolog.Ctx(context.Background()).With().Fields(map[string]any{
					"schedule_id": uuid.New().String(),
				}).Logger().WithContext(jobCtx)

				if err := s.SendEmailJob(jobCtx); err != nil {
					zerolog.Ctx(ctx).Error().Fields(map[string]any{
						"error": err,
					}).Msg("email schedule send email error")
				}
				jobCancel()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *EmailSchedule) SendEmailJob(ctx context.Context) error {
	emailJobs, err := s.repo.GrabEmailJob(ctx)
	if err != nil {
		return fmt.Errorf("repo.GrabEmailJob error: %w", err)
	}

	for _, emailJob := range emailJobs {
		// get email sender by sender_id
		emailSender, err := s.repo.GetEmailSenderByID(ctx, emailJob.SenderID)
		if err != nil {
			zerolog.Ctx(ctx).Error().Fields(map[string]any{
				"email_job_id": emailJob.ID,
				"error":        err,
			}).Msg("repo.GetEmailSenderByID error")

			continue
		}

		sentEmailsCount, err := s.repo.CountSentEmailsLast24Hours(ctx, emailSender.ID)
		if err != nil {
			zerolog.Ctx(ctx).Error().Fields(map[string]any{
				"email_job_id": emailJob.ID,
				"error":        err,
			}).Msg("repo.CountSentEmailsLast24Hours error")

			continue
		}

		remainBudget := int64(float64(emailSender.RateLimit)*0.9) - sentEmailsCount

		for i := 1; i <= defaultEmailCountPerMinute; i++ {
			if remainBudget <= 0 {
				zerolog.Ctx(ctx).Info().Fields(map[string]any{
					"sender_id":         emailSender.ID,
					"sent_emails_count": sentEmailsCount,
					"rate_limit":        emailSender.RateLimit,
				}).Msg("email sender rate limit exceeded")

				break
			}

			if err := s.executeJob(ctx, emailJob, emailSender); err != nil {
				zerolog.Ctx(ctx).Error().Fields(map[string]any{
					"email_job_id": emailJob.ID,
					"email_sander": emailSender.Name,
					"error":        err,
				}).Msg("executeJob error")

				break
			}

			remainBudget--
		}
	}

	return nil
}

func (s *EmailSchedule) executeJob(
	ctx context.Context,
	emailJob *entities.EmailJob,
	emailSender *entities.EmailSender,
) error {
	err := s.repo.WithTx(ctx, func(ctx context.Context) error {
		emailJobEntity, err := s.repo.GetEmailJobByIDForUpdate(ctx, emailJob.ID)
		if err != nil {
			return fmt.Errorf("repo.GetEmailJobByIDForUpdate error: %w", err)
		}

		// skip if email job status is not pending or processing
		if emailJobEntity.Status != email.EmailJobStatusPending && emailJobEntity.Status != email.EmailJobStatusProcessing {
			return nil
		}

		if emailJobEntity.Status == email.EmailJobStatusPending {
			if err := s.repo.UpdateEmailJobStats(ctx, emailJobEntity.ID, email.EmailJobStatusProcessing); err != nil {
				return fmt.Errorf("repo.UpdateEmailJobStats error: %w", err)
			}
		}

		// get email receiver from email log
		emailLog, err := s.repo.GrabPendingEmailLogByJobID(ctx, emailJobEntity.ID)
		if err != nil {
			if errors.Is(err, commonErrors.ErrDataNotFound) {
				return nil
			}

			return fmt.Errorf("repo.GrabPendingEmailLogByJobID error: %w", err)
		}

		pendingEmailCount, err := s.repo.CountPendingEmailLogsByJobID(ctx, emailJobEntity.ID)
		if err != nil {
			return fmt.Errorf("repo.CountPendingEmailLogsByJobID error: %w", err)
		}

		memo := ""
		emailLogStatus := email.EmailLogStatusSuccess

		// send email
		var emailPayload domain.SendEmailParams
		err = json.Unmarshal([]byte(emailJobEntity.Payload), &emailPayload)
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
			Images:      emailPayload.Images,
			SenderEmail: emailSender.Email,
			SenderPwd:   emailSender.Key,
		})
		if err != nil {
			zerolog.Ctx(ctx).Error().Fields(map[string]any{
				"error":        err,
				"email_log_id": emailLog.ID,
			}).Msg("emailRepo.SendEmail error")

			memo = fmt.Sprintf("failed to send email: %s", err.Error())
			emailLogStatus = email.EmailLogStatusFailed
		}

		// update email log information
		err = s.repo.UpdateEmailLog(ctx, domain.UpdateEmailLogParam{
			ID:     emailLog.ID,
			Status: &emailLogStatus,
			Memo:   memo,
		})
		if err != nil {
			return fmt.Errorf("repo.UpdateEmailLog error: %w", err)
		}

		updateEmailJobParam := domain.UpdateEmailJobParam{
			JobID: emailJobEntity.ID,
		}

		if emailLogStatus == email.EmailLogStatusSuccess {
			updateEmailJobParam.IncreaseSuccessCount = 1
			emailJobEntity.SuccessCount++
		}

		if pendingEmailCount -1 <= 0 {
			updateEmailJobParam.Status = email.EmailJobStatusPartiallySuccess.ToPointer()
			if emailJobEntity.SuccessCount >= emailJobEntity.ExpectedReciverCount {
				updateEmailJobParam.Status = email.EmailJobStatusSuccess.ToPointer()
			}

			if emailJobEntity.SuccessCount == 0 {
				updateEmailJobParam.Status = email.EmailJobStatusFailed.ToPointer()
			}
		}

		// update email job information
		err = s.repo.UpdateEmailJob(ctx, updateEmailJobParam)
		if err != nil {
			return fmt.Errorf("repo.UpdateEmailJob error: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("repo.WithTx error: %w", err)
	}

	return nil
}
