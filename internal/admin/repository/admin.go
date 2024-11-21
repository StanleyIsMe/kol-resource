package repository

import (
	"context"
	"database/sql"
	"errors"
	"kolresource/internal/admin/domain"
	"kolresource/internal/admin/domain/entities"
	model "kolresource/internal/db/sqlboiler"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAdminByUserName(ctx context.Context, username string) (*entities.Admin, error) {
	adminModel, err := model.Admins(qm.Where("username = ?", username)).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, domain.QueryRecordError{Err: err}
	}

	return r.newAdminFromModel(adminModel)
}

func (r *AdminRepository) CreateAdmin(ctx context.Context, adminEntity *entities.Admin) (*entities.Admin, error) {
	adminUUID, err := uuid.NewV7()
	if err != nil {
		return nil, domain.GenerateUUIDError{Err: err}
	}

	adminModel := &model.Admin{
		ID:       adminUUID.String(),
		Name:     adminEntity.Name,
		Username: adminEntity.Username,
		Password: adminEntity.Password,
		Salt:     adminEntity.Salt,
	}

	err = adminModel.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, domain.InsertRecordError{Err: err}
	}

	return r.newAdminFromModel(adminModel)
}

func (r *AdminRepository) newAdminFromModel(adminModel *model.Admin) (*entities.Admin, error) {
	adminUUID, err := uuid.Parse(adminModel.ID)
	if err != nil {
		return nil, domain.GenerateUUIDError{Err: err}
	}

	return &entities.Admin{
		ID:        adminUUID,
		Name:      adminModel.Name,
		Username:  adminModel.Username,
		Password:  adminModel.Password,
		Salt:      adminModel.Salt,
		CreatedAt: adminModel.CreatedAt,
		UpdatedAt: adminModel.UpdatedAt,
	}, nil
}
