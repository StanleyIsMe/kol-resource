package domain

import (
	"context"
	"kolresource/internal/email"
	"kolresource/internal/email/domain/entities"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error

	CreateEmailSender(ctx context.Context, sender *entities.EmailSender) error
	GetEmailSenderByID(ctx context.Context, id int64) (*entities.EmailSender, error)
	GetEmailSenderByEmail(ctx context.Context, email string) (*entities.EmailSender, error)
	AllEmailSenders(ctx context.Context) ([]*entities.EmailSender, error)

	CreateEmailJob(ctx context.Context, job *entities.EmailJob) (*entities.EmailJob, error)
	UpdateEmailJob(ctx context.Context, job *entities.EmailJob) error
	GrabEmailJob(ctx context.Context, status email.EmailJobStatus) (*entities.EmailJob, error)
	ListEmailJobs(ctx context.Context, params *ListEmailJobsParams) ([]*entities.EmailJob, error)

	BatchCreateEmailLogs(ctx context.Context, logs []*entities.EmailLog) error
	UpdateEmailLog(ctx context.Context, log *entities.EmailLog) error
	GetEmailLog(ctx context.Context, id int64) (*entities.EmailLog, error)
	ListEmailLogs(ctx context.Context, params *ListEmailLogsParams) ([]*entities.EmailLog, error)
}

type GrabEmailJobParams struct {
	ID int64
}

type ListEmailLogsParams struct {
	JobID  int64
	Status *email.EmailLogStatus
}

type ListEmailJobsParams struct {
	Status *email.EmailJobStatus
	Page   int64
	Size   int64
}
