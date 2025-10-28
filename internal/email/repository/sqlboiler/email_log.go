package sqlboiler

import (
	"context"
	"fmt"
	commonErrors "kolresource/internal/common/errors"
	model "kolresource/internal/db/sqlboiler"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
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
			ProductID: emailLog.ProductID.String(),
			SenderID:  emailLog.SenderID.String(),
			MessageID: emailLog.MessageID,
			Reply:     emailLog.Reply,
			Momo:      emailLog.Memo,
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
	var qmMods []qm.QueryMod

	qmMods = append(qmMods, qm.Where("job_id = ?", params.JobID))

	if params.Status != nil {
		qmMods = append(qmMods, qm.Where("status = ?", *params.Status))
	}

	emailLogs, err := model.EmailLogs(qmMods...).All(ctx, r.db)
	if err != nil {
		return nil, commonErrors.QueryRecordError{Err: err}
	}

	emailLogEntities := make([]*entities.EmailLog, len(emailLogs))
	for index, emailLog := range emailLogs {
		emailLog := emailLog
		emailLogEntity, err := r.newEmailLogFromModel(emailLog)
		if err != nil {
			return nil, fmt.Errorf("failed to convert email log to entities: %w", err)
		}
		emailLogEntities[index] = emailLogEntity
	}

	return emailLogEntities, nil
}

func (r *EmailRepository) newEmailLogFromModel(emailLogModel *model.EmailLog) (*entities.EmailLog, error) {
	kolID, err := uuid.Parse(emailLogModel.KolID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "kol_id", UUID: emailLogModel.KolID}
	}

	productID, err := uuid.Parse(emailLogModel.ProductID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "product_id", UUID: emailLogModel.ProductID}
	}

	senderID, err := uuid.Parse(emailLogModel.SenderID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "sender_id", UUID: emailLogModel.SenderID}
	}

	return &entities.EmailLog{
		ID:        emailLogModel.ID,
		JobID:     emailLogModel.JobID,
		KolID:     kolID,
		KolName:   emailLogModel.KolName,
		Email:     emailLogModel.Email,
		Status:    email.EmailLogStatus(emailLogModel.Status),
		ProductID: productID,
		SenderID:  senderID,
		MessageID: emailLogModel.MessageID,
		Reply:     emailLogModel.Reply,
		Memo:      emailLogModel.Momo,
		CreatedAt: emailLogModel.CreatedAt,
		UpdatedAt: emailLogModel.UpdatedAt,
	}, nil
}
