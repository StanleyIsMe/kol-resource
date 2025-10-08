package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"
	kolUsecase "kolresource/internal/kol/usecase"
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
		return fmt.Errorf("email sender already exists")
	}

	emailSender := &entities.EmailSender{
		Name:       param.Name,
		Email:      param.Email,
		Key:        param.Key,
		LastSendAt: time.Now(),
	}
	emailSender.DefaultRateLimit()

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
	for _, emailSender := range emailSenders {
		emailSenders = append(emailSenders, EmailSender{
			ID:    emailSender.ID,
			Name:  emailSender.Name,
			Email: emailSender.Email,
		})
	}

	return emailSenders, nil
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
			// Payload:              param.Payload, // TODO: email content + images
			Status:        email.EmailJobStatusPending,
			LastExecuteAt: time.Now(),
		}

		emailJob, err := uc.repo.CreateEmailJob(ctx, job)
		if err != nil {
			return fmt.Errorf("repo.CreateEmailJob error: %w", err)
		}

		emailLogs := make([]*entities.EmailLog, 0, len(kols))
		for _, kol := range kols {
			emailLogs = append(emailLogs, &entities.EmailLog{
				JobID:   emailJob.ID,
				KolID:   kol.ID,
				KolName: kol.Name,
				Email:   kol.Email,
				Status:  email.EmailLogStatusPending,
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