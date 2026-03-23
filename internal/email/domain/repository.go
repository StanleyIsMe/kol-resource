package domain

import (
	"context"
	"kolresource/internal/email"
	"kolresource/internal/email/domain/entities"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=../mock/repositorymock/repository_mock.go -package=repositorymock
//nolint:interfacebloat
type Repository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error

	CreateEmailSender(ctx context.Context, sender *entities.EmailSender) error
	GetEmailSenderByID(ctx context.Context, id uuid.UUID) (*entities.EmailSender, error)
	UpdateEmailSender(ctx context.Context, param UpdateEmailSenderParam) error
	GetEmailSenderByEmail(ctx context.Context, email string) (*entities.EmailSender, error)
	AllEmailSenders(ctx context.Context) ([]*entities.EmailSender, error)

	CreateEmailJob(ctx context.Context, job *entities.EmailJob) (*entities.EmailJob, error)
	UpdateEmailJobStats(ctx context.Context, id int64, status email.JobStatus) error
	UpdateEmailJob(ctx context.Context, param UpdateEmailJobParam) error
	GetEmailJobByID(ctx context.Context, id int64) (*entities.EmailJob, error)
	GetEmailJobByIDForUpdate(ctx context.Context, id int64) (*entities.EmailJob, error)
	GrabEmailJob(ctx context.Context) ([]*entities.EmailJob, error)
	ListEmailJobs(ctx context.Context, params *ListEmailJobsParams) ([]*entities.EmailJob, int64, error)

	BatchCreateEmailLogs(ctx context.Context, logs []*entities.EmailLog) error
	UpdateEmailLog(ctx context.Context, param UpdateEmailLogParam) error
	GetEmailLog(ctx context.Context, id int64) (*entities.EmailLog, error)
	ListEmailLogs(ctx context.Context, params *ListEmailLogsParams) ([]*entities.EmailLog, error)
	GrabPendingEmailLogByJobID(ctx context.Context, jobID int64) (*entities.EmailLog, error)
	CountPendingEmailLogsByJobID(ctx context.Context, jobID int64) (int64, error)
	CountSentEmailsLast24Hours(ctx context.Context, senderID uuid.UUID) (int64, error)
}

type EmailRepository interface {
	SendEmail(ctx context.Context, param SendEmailParams) error
}

type UpdateEmailSenderParam struct {
	ID             uuid.UUID `json:"id"`
	Name           *string   `json:"name"`
	Email          *string   `json:"email"`
	Key            *string   `json:"key"`
	RateLimit      *int      `json:"rate_limit"`
	UpdatedAdminID uuid.UUID `json:"updated_admin_id"`
}

type GrabEmailJobParams struct {
	ID int64
}

type ListEmailLogsParams struct {
	JobID  int64
	Status *email.LogStatus
}

type ListEmailJobsParams struct {
	SenderEmail *string
	SenderName  *string
	ProductName *string
	Status      *email.JobStatus
	Page        int
	Size        int
}

type UpdateEmailJobParam struct {
	JobID                int64
	Status               *email.JobStatus
	IncreaseSuccessCount int
}

type UpdateEmailLogParam struct {
	ID     int64
	Status *email.LogStatus
	Memo   string
	Reply  *bool
}

type SendEmailParams struct {
	Subject     string           `json:"subject"`
	Body        string           `json:"body"`
	ToEmails    []ToEmail        `json:"-"`
	Images      []SendEmailImage `json:"images"`
	SenderName  string           `json:"-"`
	SenderEmail string           `json:"-"`
	SenderPwd   string           `json:"-"`
}

type ToEmail struct {
	Email string
	Name  string
}

type SendEmailImage struct {
	ContentID string `json:"content_id"`
	Data      string `json:"data"`
	ImageType string `json:"image_type"`
}
