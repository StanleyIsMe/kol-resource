package sqlboiler

import (
	"context"
	commonErrors "kolresource/internal/common/errors"
	model "kolresource/internal/db/sqlboiler"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (r *EmailRepository) BatchCreateEmailLogs(ctx context.Context, logs []*entities.EmailLog) error {
	tx := r.getTx(ctx)

	// due to the sqlboiler was not friendly for batch insert, we need to insert one by one
	// TODO: https://github.com/tiendc/sqlboiler-extensions
	for _, emailLog := range logs {
		emailLogModel := &model.EmailLog{
			JobID:     emailLog.JobID,
			KolID:     emailLog.KolID.String(),
			KolName:   emailLog.KolName,
			Email:     emailLog.Email,
			Status:    model.EmailLogStatus(emailLog.Status),
		}

		err := emailLogModel.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *EmailRepository) UpdateEmailLog(ctx context.Context, log *entities.EmailLog) error {
	return nil
}

func (r *EmailRepository) GetEmailLog(ctx context.Context, id int64) (*entities.EmailLog, error) {
	return nil, nil
}

func (r *EmailRepository) ListEmailLogs(ctx context.Context, params *domain.ListEmailLogsParams) ([]*entities.EmailLog, error) {
	return nil, nil
}

func (r *EmailRepository) newEmailLogFromModel(emailLogModel *model.EmailLog) (*entities.EmailLog, error) {
	kolID, err := uuid.Parse(emailLogModel.KolID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "kol_id", UUID: emailLogModel.KolID}
	}

	return &entities.EmailLog{
		ID:        emailLogModel.ID,
		JobID:     emailLogModel.JobID,
		KolID:     kolID,
		KolName:   emailLogModel.KolName,
		Email:     emailLogModel.Email,
		Status:    email.EmailLogStatus(emailLogModel.Status),
		CreatedAt: emailLogModel.CreatedAt,
		UpdatedAt: emailLogModel.UpdatedAt,
	}, nil
}
