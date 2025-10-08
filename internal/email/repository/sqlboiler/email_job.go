package sqlboiler

import (
	"context"
	"fmt"
	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"

	model "kolresource/internal/db/sqlboiler"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (r *EmailRepository) CreateEmailJob(ctx context.Context, job *entities.EmailJob) (*entities.EmailJob, error) {
	emailJobModel := &model.EmailJob{
		ID:                   job.ID,
		ExpectedReciverCount: job.ExpectedReciverCount,
		Status:               model.EmailJobStatus(job.Status),
		AdminID:              job.AdminID.String(),
		AdminName:            job.AdminName,
		SenderID:             job.SenderID.String(),
		SenderName:           job.SenderName,
		SenderEmail:          job.SenderEmail,
		ProductID:            job.ProductID.String(),
		ProductName:          job.ProductName,
		Memo:                 job.Memo,
		// Payload:             job.Payload,
		LastExecuteAt: job.LastExecuteAt,
	}

	err := emailJobModel.Insert(ctx, r.getTx(ctx), boil.Infer())
	if err != nil {
		return nil, commonErrors.InsertRecordError{Err: err}
	}

	return nil, nil
}

func (r *EmailRepository) UpdateEmailJob(ctx context.Context, job *entities.EmailJob) error {
	return nil
}

func (r *EmailRepository) GrabEmailJob(ctx context.Context, status email.EmailJobStatus) (*entities.EmailJob, error) {
	emailJobModel, err := model.EmailJobs(
		qm.Where("status = ?", model.EmailJobStatus(status)),
		qm.OrderBy("created_at ASC"),
	).One(ctx, r.db)
	if err != nil {
		return nil, commonErrors.QueryRecordError{Err: err}
	}

	return r.newEmailJobFromModel(emailJobModel)
}

func (r *EmailRepository) ListEmailJobs(ctx context.Context, params *domain.ListEmailJobsParams) ([]*entities.EmailJob, error) {
	var qmMods []qm.QueryMod

	if params.Status != nil {
		qmMods = append(qmMods, qm.Where("status = ?", model.EmailJobStatus(*params.Status)))
	}

	qmMods = append(qmMods,
		qm.OrderBy("created_at DESC"),
		qm.Offset(int(params.Page*params.Size)),
		qm.Limit(int(params.Size)),
	)

	emailJobs, err := model.EmailJobs(qmMods...).All(ctx, r.db)
	if err != nil {
		return nil, commonErrors.QueryRecordError{Err: err}
	}

	emailJobsWithKols := make([]*entities.EmailJob, len(emailJobs))
	for index, emailJob := range emailJobs {
		emailJobWithKol, err := r.newEmailJobFromModel(emailJob)
		if err != nil {
			return nil, fmt.Errorf("failed to convert email job to entities: %w", err)
		}

		emailJobsWithKols[index] = emailJobWithKol
	}

	return emailJobsWithKols, nil
}

func (r *EmailRepository) newEmailJobFromModel(emailJobModel *model.EmailJob) (*entities.EmailJob, error) {
	adminID, err := uuid.Parse(emailJobModel.AdminID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "admin_id", UUID: emailJobModel.AdminID}
	}

	senderID, err := uuid.Parse(emailJobModel.SenderID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "sender_id", UUID: emailJobModel.SenderID}
	}

	productID, err := uuid.Parse(emailJobModel.ProductID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "product_id", UUID: emailJobModel.ProductID}
	}

	return &entities.EmailJob{
		ID:                   emailJobModel.ID,
		ExpectedReciverCount: emailJobModel.ExpectedReciverCount,
		Status:               email.EmailJobStatus(emailJobModel.Status),
		AdminID:              adminID,
		AdminName:            emailJobModel.AdminName,
		SenderID:             senderID,
		SenderName:           emailJobModel.SenderName,
		SenderEmail:          emailJobModel.SenderEmail,
		ProductID:            productID,
		ProductName:          emailJobModel.ProductName,
		Memo:                 emailJobModel.Memo,
		Payload:              emailJobModel.Payload.String(),
		LastExecuteAt:        emailJobModel.LastExecuteAt,
		CreatedAt:            emailJobModel.CreatedAt,
		UpdatedAt:            emailJobModel.UpdatedAt,
	}, nil
}
