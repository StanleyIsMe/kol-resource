package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	commonErrors "kolresource/internal/common/errors"
	model "kolresource/internal/db/sqlboiler"
	"kolresource/internal/email/domain/entities"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/google/uuid"
)

func (r *EmailRepository) AllEmailSenders(ctx context.Context) ([]*entities.EmailSender, error) {
	emailSendersModel, err := model.EmailSenders().All(ctx, r.db)
	if err != nil {
		return nil, commonErrors.QueryRecordError{Err: err}
	}

	emailSenders := make([]*entities.EmailSender, 0, len(emailSendersModel))
	for _, emailSender := range emailSendersModel {
		emailSenderEntity, err := r.newEmailSenderFromModel(emailSender)
		if err != nil {
			return nil, fmt.Errorf("failed to convert email sender model to entity: %w", err)
		}

		emailSenders = append(emailSenders, emailSenderEntity)
	}

	return emailSenders, nil
}

func (r *EmailRepository) CreateEmailSender(ctx context.Context, sender *entities.EmailSender) error {
	emailSenderModel := &model.EmailSender{
		ID:         sender.ID.String(),
		Name:       sender.Name,
		Email:      sender.Email,
		Key:        sender.Key,
		RateLimit:  sender.RateLimit,
		LastSendAt: sender.LastSendAt,
	}

	err := emailSenderModel.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return commonErrors.InsertRecordError{Err: err}
	}

	return nil
}

func (r *EmailRepository) GetEmailSenderByID(ctx context.Context, id int64) (*entities.EmailSender, error) {
	emailSenderModel, err := model.EmailSenders(qm.Where("id = ?", id)).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonErrors.ErrDataNotFound
		}

		return nil, commonErrors.QueryRecordError{Err: err}
	}

	return r.newEmailSenderFromModel(emailSenderModel)
}

func (r *EmailRepository) GetEmailSenderByEmail(ctx context.Context, email string) (*entities.EmailSender, error) {
	emailSenderModel, err := model.EmailSenders(qm.Where("email = ?", email)).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonErrors.ErrDataNotFound
		}

		return nil, commonErrors.QueryRecordError{Err: err}
	}

	return r.newEmailSenderFromModel(emailSenderModel)
}

func (r *EmailRepository) newEmailSenderFromModel(emailSenderModel *model.EmailSender) (*entities.EmailSender, error) {
	emailSenderUUID, err := uuid.Parse(emailSenderModel.ID)
	if err != nil {
		return nil, commonErrors.UUIDInvalidError{Field: "id", UUID: emailSenderModel.ID}
	}

	return &entities.EmailSender{
		ID:         emailSenderUUID,
		Name:       emailSenderModel.Name,
		Email:      emailSenderModel.Email,
		Key:        emailSenderModel.Key,
		RateLimit:  emailSenderModel.RateLimit,
		LastSendAt: emailSenderModel.LastSendAt,
		CreatedAt:  emailSenderModel.CreatedAt,
		UpdatedAt:  emailSenderModel.UpdatedAt,
	}, nil
}
