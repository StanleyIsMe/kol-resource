package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"
	kolUsecase "kolresource/internal/kol/usecase"

	"github.com/google/uuid"
)

type EmailUseCaseImpl struct {
	repo       domain.Repository
	kolUsecase kolUsecase.KolUseCase
}

var _ EmailUseCase = (*EmailUseCaseImpl)(nil)

func NewEmailUseCaseImpl(repo domain.Repository, kolUsecase kolUsecase.KolUseCase) *EmailUseCaseImpl {
	return &EmailUseCaseImpl{repo: repo, kolUsecase: kolUsecase}
}

func (uc *EmailUseCaseImpl) CreateEmailSender(ctx context.Context, param CreateEmailSenderParam) error {
	existEmailSender, err := uc.repo.GetEmailSenderByEmail(ctx, param.Email)
	if err != nil && !errors.Is(err, commonErrors.ErrDataNotFound) {
		return fmt.Errorf("repo.GetEmailSenderByEmail error: %w", err)
	}

	if existEmailSender != nil {
		return commonErrors.DuplicatedResourceError{
			Resource: "email sender",
			Name:     param.Email,
		}
	}

	emailSender := &entities.EmailSender{
		Name:           param.Name,
		Email:          param.Email,
		Key:            param.Key,
		RateLimit:      param.RateLimit,
		UpdatedAdminID: param.UpdatedAdminID,
		LastSendAt:     time.Now(),
	}
	emailSender.SetDailyRateLimit(param.RateLimit)

	err = uc.repo.CreateEmailSender(ctx, emailSender)
	if err != nil {
		return fmt.Errorf("repo.CreateEmailSender error: %w", err)
	}

	return nil
}

func (uc *EmailUseCaseImpl) ListEmailSenders(ctx context.Context) ([]EmailSender, error) {
	entitiesEmailSenders, err := uc.repo.AllEmailSenders(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo.ListEmailSenders error: %w", err)
	}

	emailSenders := make([]EmailSender, 0, len(entitiesEmailSenders))
	for _, emailSender := range entitiesEmailSenders {
		emailSenders = append(emailSenders, EmailSender{
			ID:         emailSender.ID,
			Name:       emailSender.Name,
			Email:      emailSender.Email,
			RateLimit:  emailSender.RateLimit,
			LastSendAt: emailSender.LastSendAt,
		})
	}

	return emailSenders, nil
}

func (uc *EmailUseCaseImpl) UpdateEmailSender(ctx context.Context, param UpdateEmailSenderParam) error {
	err := uc.repo.UpdateEmailSender(ctx, domain.UpdateEmailSenderParam{
		ID:             param.ID,
		Name:           param.Name,
		Email:          param.Email,
		Key:            param.Key,
		RateLimit:      param.RateLimit,
		UpdatedAdminID: param.UpdatedAdminID,
	})
	if err != nil {
		return fmt.Errorf("repo.UpdateEmailSender error: %w", err)
	}

	return nil
}

func (uc *EmailUseCaseImpl) GetEmailSender(ctx context.Context, id uuid.UUID) (*EmailSender, error) {
	emailSender, err := uc.repo.GetEmailSenderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo.GetEmailSenderByID error: %w", err)
	}

	return &EmailSender{
		ID:         emailSender.ID,
		Name:       emailSender.Name,
		Email:      emailSender.Email,
		RateLimit:  emailSender.RateLimit,
		LastSendAt: emailSender.LastSendAt,
	}, nil
}

func (uc *EmailUseCaseImpl) SendEmail(ctx context.Context, param SendEmailParam) error {
	emailSender, err := uc.repo.GetEmailSenderByID(ctx, param.SenderID)
	if err != nil {
		if errors.Is(err, commonErrors.ErrDataNotFound) {
			return commonErrors.NotFoundError{
				Resource: "email sender",
				ID:       param.SenderID,
			}
		}

		return fmt.Errorf("repo.GetEmailSenderByID error: %w", err)
	}

	product, err := uc.kolUsecase.GetProductByID(ctx, param.ProductID)
	if err != nil {
		if errors.Is(err, commonErrors.ErrDataNotFound) {
			return commonErrors.NotFoundError{
				Resource: "product",
				ID:       param.ProductID,
			}
		}

		return fmt.Errorf("kolUsecase.GetProductByID error: %w", err)
	}

	kols, err := uc.kolUsecase.ListKolEmailsByIDs(ctx, param.KolIDs)
	if err != nil {
		return fmt.Errorf("kolUsecase.ListKolEmailsByIDs error: %w", err)
	}

	if len(kols) == 0 {
		return commonErrors.NotFoundError{
			Resource: "kols",
			ID:       param.KolIDs,
		}
	}

	err = uc.repo.WithTx(ctx, func(ctx context.Context) error {
		job := &entities.EmailJob{
			ExpectedReciverCount: len(kols),
			AdminID:              param.UpdatedAdminID,
			AdminName:            param.UpdatedAdminName,
			SenderID:             emailSender.ID,
			SenderName:           emailSender.Name,
			SenderEmail:          emailSender.Email,
			ProductID:            param.ProductID,
			ProductName:          product.Name,
			UpdatedAdminID:       param.UpdatedAdminID,
			Status:               email.EmailJobStatusPending,
			LastExecuteAt:        time.Now(),
		}

		payload := domain.SendEmailParams{
			Subject: param.Subject,
			Body:    param.EmailContent,
			Images:  make([]domain.SendEmailImage, 0, len(param.Images)),
		}
		for _, image := range param.Images {
			payload.Images = append(payload.Images, domain.SendEmailImage{
				ContentID: image.ContentID,
				Data:      image.Data,
				ImageType: image.ImageType,
			})
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("json.Marshal error: %w", err)
		}
		job.Payload = string(payloadBytes)

		emailJob, err := uc.repo.CreateEmailJob(ctx, job)
		if err != nil {
			return fmt.Errorf("repo.CreateEmailJob error: %w", err)
		}

		emailLogs := make([]*entities.EmailLog, 0, len(kols))
		for _, kol := range kols {
			kol := kol
			emailLogs = append(emailLogs, &entities.EmailLog{
				JobID:     emailJob.ID,
				KolID:     kol.ID,
				KolName:   kol.Name,
				Email:     kol.Email,
				Status:    email.EmailLogStatusPending,
				ProductID: param.ProductID,
				SenderID:  emailSender.ID,
				MessageID: "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		err = uc.repo.BatchCreateEmailLogs(ctx, emailLogs)
		if err != nil {
			return fmt.Errorf("repo.BatchCreateEmailLogs error: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("repo.WithTx error: %w", err)
	}

	return nil
}

func (uc *EmailUseCaseImpl) ListEmailJobs(ctx context.Context, param ListEmailJobsParam) (*ListEmailJobsResponse, error) {
	emailJobs, total, err := uc.repo.ListEmailJobs(ctx, &domain.ListEmailJobsParams{
		SenderEmail: param.SenderEmail,
		SenderName:  param.SenderName,
		ProductName: param.ProductName,
		Status:   param.Status,
		Page:     param.Page,
		Size:     param.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.ListEmailJobs error: %w", err)
	}

	emailJobsResponse := &ListEmailJobsResponse{
		EmailJobs: make([]EmailJob, 0, len(emailJobs)),
		Total:     total,
	}

	for _, emailJob := range emailJobs {
		emailJob := emailJob
		emailJobsResponse.EmailJobs = append(emailJobsResponse.EmailJobs, EmailJob{
			ID:                   emailJob.ID,
			ExpectedReciverCount: emailJob.ExpectedReciverCount,
			SuccessCount:         emailJob.SuccessCount,
			SenderID:             emailJob.SenderID,
			SenderName:           emailJob.SenderName,
			SenderEmail:          emailJob.SenderEmail,
			AdminID:              emailJob.AdminID,
			AdminName:            emailJob.AdminName,
			ProductID:            emailJob.ProductID,
			ProductName:          emailJob.ProductName,
			Memo:                 emailJob.Memo,
			Status:               emailJob.Status,
			CreatedAt:            emailJob.CreatedAt,
			UpdatedAt:            emailJob.UpdatedAt,
			LastExecuteAt:        emailJob.LastExecuteAt,
		})
	}

	return emailJobsResponse, nil
}

func (uc *EmailUseCaseImpl) GetEmailJob(ctx context.Context, id int64) (*EmailJob, error) {
	emailJob, err := uc.repo.GetEmailJobByID(ctx, id)
	if err != nil {
		if errors.Is(err, commonErrors.ErrDataNotFound) {
			return nil, commonErrors.NotFoundError{
				Resource: "email_job",
				ID:       id,
			}
		}

		return nil, fmt.Errorf("repo.GetEmailJobByID error: %w", err)
	}

	return &EmailJob{
		ID:                   emailJob.ID,
		ExpectedReciverCount: emailJob.ExpectedReciverCount,
		SuccessCount:         emailJob.SuccessCount,
		SenderID:             emailJob.SenderID,
		SenderName:           emailJob.SenderName,
		SenderEmail:          emailJob.SenderEmail,
		AdminID:              emailJob.AdminID,
		AdminName:            emailJob.AdminName,
		ProductID:            emailJob.ProductID,
		ProductName:          emailJob.ProductName,
		Memo:                 emailJob.Memo,
		Status:               emailJob.Status,
		CreatedAt:            emailJob.CreatedAt,
		UpdatedAt:            emailJob.UpdatedAt,
		LastExecuteAt:        emailJob.LastExecuteAt,
	}, nil
}

func (uc *EmailUseCaseImpl) ListEmailLogs(ctx context.Context, param ListEmailLogsParam) ([]EmailLog, error) {
	emailLogEntities, err := uc.repo.ListEmailLogs(ctx, &domain.ListEmailLogsParams{
		JobID:  param.JobID,
		Status: param.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.ListEmailLogs error: %w", err)
	}

	emailLogs := make([]EmailLog, 0, len(emailLogEntities))
	for _, emailLog := range emailLogEntities {
		emailLog := emailLog
		emailLogs = append(emailLogs, EmailLog{
			ID:        emailLog.ID,
			Email:     emailLog.Email,
			Reply:     emailLog.Reply,
			KolID:     emailLog.KolID,
			KolName:   emailLog.KolName,
			Status:    emailLog.Status,
			Memo:      emailLog.Memo,
			SendedAt:  emailLog.SendedAt,
		})
	}

	return emailLogs, nil
}

func (uc *EmailUseCaseImpl) CancelEmailJob(ctx context.Context, id int64) error {
	emailJob, err := uc.repo.GetEmailJobByID(ctx, id)
	if err != nil {
		if errors.Is(err, commonErrors.ErrDataNotFound) {
			return commonErrors.NotFoundError{
				Resource: "email_job",
				ID:       id,
			}
		}

		return fmt.Errorf("repo.GetEmailJobByID error: %w", err)
	}

	if !emailJob.Status.CanCancel() {
		return fmt.Errorf("email job status is not cancelable")
	}

	err = uc.repo.WithTx(ctx, func(ctx context.Context) error {
		err = uc.repo.UpdateEmailJobStats(ctx, id, email.EmailJobStatusCanceled)
		if err != nil {
			return fmt.Errorf("repo.UpdateEmailJobStats error: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("repo.WithTx error: %w", err)
	}

	return nil
}

func (uc *EmailUseCaseImpl) StartEmailJob(ctx context.Context, id int64) error {
	emailJob, err := uc.repo.GetEmailJobByID(ctx, id)
	if err != nil {
		if errors.Is(err, commonErrors.ErrDataNotFound) {
			return commonErrors.NotFoundError{
				Resource: "email_job",
				ID:       id,
			}
		}

		return fmt.Errorf("repo.GetEmailJobByID error: %w", err)
	}

	if !emailJob.Status.CanStart() {
		return fmt.Errorf("email job status is not startable")
	}

	err = uc.repo.WithTx(ctx, func(ctx context.Context) error {
		err = uc.repo.UpdateEmailJobStats(ctx, id, email.EmailJobStatusPending)
		if err != nil {
			return fmt.Errorf("repo.UpdateEmailJobStats error: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("repo.WithTx error: %w", err)
	}

	return nil
}
