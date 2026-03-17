package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"
	"strings"
	"time"

	model "kolresource/internal/db/sqlboiler"

	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/google/uuid"

	"github.com/aarondl/sqlboiler/v4/boil"
)

func (r *EmailRepository) CreateEmailJob(ctx context.Context, job *entities.EmailJob) (*entities.EmailJob, error) {
	emailJobModel := &model.EmailJob{
		ExpectedReciverCount: job.ExpectedReciverCount,
		Status:               model.EmailJobStatus(job.Status),
		AdminID:              job.AdminID.String(),
		AdminName:            job.AdminName,
		SenderID:             job.SenderID.String(),
		SenderName:           job.SenderName,
		SenderEmail:          job.SenderEmail,
		ProductID:            job.ProductID.String(),
		ProductName:          job.ProductName,
		UpdatedAdminID:       job.UpdatedAdminID.String(),
		Memo:                 job.Memo,
		Payload:              types.JSON(job.Payload),
		LastExecuteAt:        job.LastExecuteAt,
	}

	err := emailJobModel.Insert(ctx, r.getTx(ctx), boil.Infer())
	if err != nil {
		return nil, commonErrors.InsertRecordError{Err: err}
	}

	return r.newEmailJobFromModel(emailJobModel)
}

func (r *EmailRepository) UpdateEmailJobStats(ctx context.Context, id int64, status email.EmailJobStatus) error {
	emailJobModel, err := model.EmailJobs(
		qm.Where("id = ?", id),
		qm.For("UPDATE"),
	).One(ctx, r.getTx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commonErrors.ErrDataNotFound
		}

		return commonErrors.QueryRecordError{Err: err}
	}

	emailJobModel.Status = model.EmailJobStatus(status)

	_, err = emailJobModel.Update(ctx, r.getTx(ctx), boil.Infer())
	if err != nil {
		return commonErrors.UpdateRecordError{Err: err}
	}

	return nil
}

func (r *EmailRepository) UpdateEmailJob(ctx context.Context, param domain.UpdateEmailJobParam) error {
	var (
		setParts []string
		args     []interface{}
	)

	argIndex := 1

	if param.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, model.EmailJobStatus(*param.Status))
		argIndex++
	}

	if param.IncreaseSuccessCount > 0 {
		setParts = append(setParts, fmt.Sprintf("success_count = success_count + $%d", argIndex))
		args = append(args, param.IncreaseSuccessCount)
		argIndex++
	}

	setParts = append(setParts, fmt.Sprintf("last_execute_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, param.JobID)

	query := fmt.Sprintf(
		`UPDATE email_job 
		 SET %s 
		 WHERE id = $%d`,
		strings.Join(setParts, ", "),
		argIndex,
	)

	result, err := queries.Raw(query, args...).ExecContext(ctx, r.getTx(ctx))
	if err != nil {
		return commonErrors.UpdateRecordError{Err: err}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return commonErrors.UpdateRecordError{Err: err}
	}

	if rowsAffected == 0 {
		return commonErrors.ErrDataNotFound
	}

	return nil
}

func (r *EmailRepository) GetEmailJobByID(ctx context.Context, id int64) (*entities.EmailJob, error) {
	emailJobModel, err := model.EmailJobs(
		qm.Where("id = ?", id),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonErrors.ErrDataNotFound
		}

		return nil, commonErrors.QueryRecordError{Err: err}
	}

	return r.newEmailJobFromModel(emailJobModel)
}

func (r *EmailRepository) GetEmailJobByIDForUpdate(ctx context.Context, id int64) (*entities.EmailJob, error) {
	emailJobModel, err := model.EmailJobs(
		qm.Where("id = ?", id),
		qm.For("UPDATE"),
	).One(ctx, r.getTx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonErrors.ErrDataNotFound
		}

		return nil, commonErrors.QueryRecordError{Err: err}
	}

	return r.newEmailJobFromModel(emailJobModel)
}

// GrabEmailJob grabs the email job with the status of pending or processing
func (r *EmailRepository) GrabEmailJob(ctx context.Context) ([]*entities.EmailJob, error) {
	query := `
		WITH ranked AS (
			SELECT *,
				ROW_NUMBER() OVER (PARTITION BY sender_id ORDER BY created_at ASC) AS rn
			FROM email_job
			WHERE status IN ('pending', 'processing')
		)
		SELECT * 
		FROM ranked
		WHERE rn = 1
		FOR UPDATE SKIP LOCKED;
	`
	
	var emailJobModels []*model.EmailJob
	err := queries.Raw(query).Bind(ctx, r.db, &emailJobModels)
	if err != nil {
		return nil, commonErrors.QueryRecordError{Err: err}
	}

	emailJobs := make([]*entities.EmailJob, len(emailJobModels))
	for index, emailJobModel := range emailJobModels {
		emailJob, err := r.newEmailJobFromModel(emailJobModel)
		if err != nil {
			return nil, fmt.Errorf("failed to convert email job to entities: %w", err)
		}
		emailJobs[index] = emailJob
	}

	return emailJobs, nil
}

func (r *EmailRepository) ListEmailJobs(ctx context.Context, params *domain.ListEmailJobsParams) ([]*entities.EmailJob, int64, error) {
	var qmMods []qm.QueryMod

	if params.Status != nil {
		qmMods = append(qmMods, qm.Where("status = ?", model.EmailJobStatus(*params.Status)))
	}

	if params.SenderEmail != nil {
		qmMods = append(qmMods, qm.Where("sender_email LIKE ?", "%"+*params.SenderEmail+"%"))
	}

	if params.SenderName != nil {
		qmMods = append(qmMods, qm.Where("sender_name LIKE ?", "%"+*params.SenderName+"%"))
	}

	if params.ProductName != nil {
		qmMods = append(qmMods, qm.Where("product_name LIKE ?", "%"+*params.ProductName+"%"))
	}

	count, err := model.EmailJobs(qmMods...).Count(ctx, r.db)
	if err != nil {
		return nil, 0, commonErrors.QueryRecordError{Err: err}
	}

	qmMods = append(qmMods,
		qm.OrderBy("updated_at DESC"),
		qm.Offset((params.Page-1)*params.Size),
		qm.Limit(params.Size),
	)

	emailJobs, err := model.EmailJobs(qmMods...).All(ctx, r.db)
	if err != nil {
		return nil, 0, commonErrors.QueryRecordError{Err: err}
	}

	emailJobsWithKols := make([]*entities.EmailJob, len(emailJobs))
	for index, emailJob := range emailJobs {
		emailJobWithKol, err := r.newEmailJobFromModel(emailJob)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert email job to entities: %w", err)
		}

		emailJobsWithKols[index] = emailJobWithKol
	}

	return emailJobsWithKols, count, nil
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

	updatedAdminID, err := uuid.Parse(emailJobModel.UpdatedAdminID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "updated_admin_id", UUID: emailJobModel.UpdatedAdminID}
	}

	return &entities.EmailJob{
		ID:                   emailJobModel.ID,
		ExpectedReciverCount: emailJobModel.ExpectedReciverCount,
		SuccessCount:         emailJobModel.SuccessCount,
		Status:               email.EmailJobStatus(emailJobModel.Status),
		AdminID:              adminID,
		AdminName:            emailJobModel.AdminName,
		SenderID:             senderID,
		SenderName:           emailJobModel.SenderName,
		SenderEmail:          emailJobModel.SenderEmail,
		UpdatedAdminID:       updatedAdminID,
		ProductID:            productID,
		ProductName:          emailJobModel.ProductName,
		Memo:                 emailJobModel.Memo,
		Payload:              emailJobModel.Payload.String(),
		LastExecuteAt:        emailJobModel.LastExecuteAt,
		CreatedAt:            emailJobModel.CreatedAt,
		UpdatedAt:            emailJobModel.UpdatedAt,
	}, nil
}
