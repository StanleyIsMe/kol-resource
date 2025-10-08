package usecase

import (
	"context"

	"github.com/google/uuid"
)

type EmailUseCase interface {
	CreateEmailSender(ctx context.Context, param CreateEmailSenderParam) error
	ListEmailSenders(ctx context.Context) ([]EmailSender, error)
	SendEmail(ctx context.Context, param SendEmailParam) error
}

type CreateEmailSenderParam struct {
	UpdatedAdminID   uuid.UUID `json:"updated_admin_id"`
	UpdatedAdminName string    `json:"updated_admin_name"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Key              string    `json:"key"`
}

type EmailSender struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type SendEmailParam struct {
	Subject          string           `json:"subject"`
	EmailContent     string           `json:"email_content"`
	KolIDs           []uuid.UUID      `json:"kol_ids"`
	ProductID        uuid.UUID        `json:"product_id"`
	UpdatedAdminID   uuid.UUID        `json:"updated_admin_id"`
	UpdatedAdminName string           `json:"updated_admin_name"`
	Images           []SendEmailImage `json:"images"`
	SenderID         int64            `json:"sender_id"`
}

type SendEmailImage struct {
	ContentID string
	Data      string
	ImageType string
}
