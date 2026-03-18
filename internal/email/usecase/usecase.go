package usecase

import (
	"context"
	"kolresource/internal/email"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./usecase.go -destination=../mock/usecasemock/usecase_mock.go -package=usecasemock
type EmailUseCase interface {
	CreateEmailSender(ctx context.Context, param CreateEmailSenderParam) error
	ListEmailSenders(ctx context.Context) ([]EmailSender, error)
	UpdateEmailSender(ctx context.Context, param UpdateEmailSenderParam) error
	GetEmailSender(ctx context.Context, id uuid.UUID) (*EmailSender, error)
	SendEmail(ctx context.Context, param SendEmailParam) error
	ListEmailJobs(ctx context.Context, param ListEmailJobsParam) (*ListEmailJobsResponse, error)
	GetEmailJob(ctx context.Context, id int64) (*EmailJob, error)
	ListEmailLogs(ctx context.Context, param ListEmailLogsParam) ([]EmailLog, error)
	CancelEmailJob(ctx context.Context, id int64) error
	StartEmailJob(ctx context.Context, id int64) error
}

type CreateEmailSenderParam struct {
	UpdatedAdminID   uuid.UUID `json:"updated_admin_id"`
	UpdatedAdminName string    `json:"updated_admin_name"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Key              string    `json:"key"`
	RateLimit        int       `json:"rate_limit"`
}

type UpdateEmailSenderParam struct {
	ID             uuid.UUID `json:"id"`
	Name           *string   `json:"name"`
	Email          *string   `json:"email"`
	Key            *string   `json:"key"`
	RateLimit      *int      `json:"rate_limit"`
	UpdatedAdminID uuid.UUID `json:"updated_admin_id"`
}

type EmailSender struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	RateLimit  int       `json:"rate_limit"`
	LastSendAt time.Time `json:"last_send_at"`
}

type SendEmailParam struct {
	Subject          string           `json:"subject"`
	EmailContent     string           `json:"email_content"`
	KolIDs           []uuid.UUID      `json:"kol_ids"`
	ProductID        uuid.UUID        `json:"product_id"`
	UpdatedAdminID   uuid.UUID        `json:"updated_admin_id"`
	UpdatedAdminName string           `json:"updated_admin_name"`
	Images           []SendEmailImage `json:"images"`
	SenderID         uuid.UUID        `json:"sender_id"`
}

type SendEmailImage struct {
	ContentID string
	Data      string
	ImageType string
}

type ListEmailJobsParam struct {
	SenderEmail *string
	SenderName  *string
	ProductName *string
	Status      *email.EmailJobStatus
	Page        int
	PageSize    int
}

type ListEmailJobsResponse struct {
	EmailJobs []EmailJob `json:"email_jobs"`
	Total     int64      `json:"total"`
}

type EmailJob struct {
	ID                   int64                `json:"id"`
	ExpectedReciverCount int                  `json:"expected_reciver_count"`
	SuccessCount         int                  `json:"success_count"`
	SenderID             uuid.UUID            `json:"sender_id"`
	SenderName           string               `json:"sender_name"`
	SenderEmail          string               `json:"sender_email"`
	AdminID              uuid.UUID            `json:"admin_id"`
	AdminName            string               `json:"admin_name"`
	ProductID            uuid.UUID            `json:"product_id"`
	ProductName          string               `json:"product_name"`
	Memo                 string               `json:"memo"`
	Status               email.EmailJobStatus `json:"status"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	LastExecuteAt        time.Time            `json:"last_execute_at"`
}

type ListEmailLogsParam struct {
	JobID  int64 `json:"job_id"`
	Status *email.EmailLogStatus
}

type EmailLog struct {
	ID        int64                `json:"id"`
	Email     string               `json:"email"`
	Reply     bool                 `json:"reply"`
	KolID     uuid.UUID            `json:"kol_id"`
	KolName   string               `json:"kol_name"`
	Status    email.EmailLogStatus `json:"status"`
	Memo      string               `json:"memo"`
	SendedAt  time.Time            `json:"sended_at"`
}
