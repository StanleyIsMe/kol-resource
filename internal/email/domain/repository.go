package domain

import (
	"context"
	"kolresource/internal/email"
	"kolresource/internal/email/domain/entities"

	"github.com/google/uuid"
)

//nolint:interfacebloat
type Repository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error

	CreateEmailSender(ctx context.Context, sender *entities.EmailSender) error
	GetEmailSenderByID(ctx context.Context, id uuid.UUID) (*entities.EmailSender, error)
	UpdateEmailSender(ctx context.Context, param UpdateEmailSenderParam) error
	GetEmailSenderByEmail(ctx context.Context, email string) (*entities.EmailSender, error)
	AllEmailSenders(ctx context.Context) ([]*entities.EmailSender, error)

	CreateEmailJob(ctx context.Context, job *entities.EmailJob) (*entities.EmailJob, error)
	UpdateEmailJobStats(ctx context.Context, id int64, status email.EmailJobStatus) error
	GetEmailJobByID(ctx context.Context, id int64) (*entities.EmailJob, error)
	GrabEmailJob(ctx context.Context, status email.EmailJobStatus) (*entities.EmailJob, error)
	ListEmailJobs(ctx context.Context, params *ListEmailJobsParams) ([]*entities.EmailJob, int64, error)

	BatchCreateEmailLogs(ctx context.Context, logs []*entities.EmailLog) error
	UpdateEmailLog(ctx context.Context, log *entities.EmailLog) error
	GetEmailLog(ctx context.Context, id int64) (*entities.EmailLog, error)
	ListEmailLogs(ctx context.Context, params *ListEmailLogsParams) ([]*entities.EmailLog, error)
}

type UpdateEmailSenderParam struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Key            string    `json:"key"`
	RateLimit      int       `json:"rate_limit"`
	UpdatedAdminID uuid.UUID `json:"updated_admin_id"`
}

type GrabEmailJobParams struct {
	ID int64
}

type ListEmailLogsParams struct {
	JobID  int64
	Status *email.EmailLogStatus
}

type ListEmailJobsParams struct {
	SenderID *string
	Status   *email.EmailJobStatus
	Page     int
	Size     int
}

type UpdateEmailJobParam struct {
	JobID  int64
	Status email.EmailJobStatus
}
