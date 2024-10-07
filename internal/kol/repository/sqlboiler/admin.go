package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	model "kolresource/internal/db/sqlboiler"
	"kolresource/internal/kol/domain/admin"
	"kolresource/internal/kol/entities"

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

func (r *AdminRepository) GetAdminByUserName(ctx context.Context, userName string) (*entities.Admin, error) {
	adminModel, err := model.Admins(qm.Where("user_name = ?", userName)).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, admin.ErrDataNotFound
		}

		return nil, admin.QueryRecordError{Err: err}
	}

	return r.newAdminFromModel(adminModel)
}

func (r *AdminRepository) CreateAdmin(ctx context.Context, adminEntity *entities.Admin) (*entities.Admin, error) {
	adminUUID, err := uuid.NewV7()
	if err != nil {
		return nil, admin.GenerateUUIDError{Err: err}
	}

	adminModel := &model.Admin{
		ID:       adminUUID.String(),
		Username: adminEntity.Username,
		Password: adminEntity.Password,
	}

	err = adminModel.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, admin.InsertRecordError{Err: err}
	}

	return r.newAdminFromModel(adminModel)
}

func (r *AdminRepository) newAdminFromModel(adminModel *model.Admin) (*entities.Admin, error) {
	adminUUID, err := uuid.Parse(adminModel.ID)
	if err != nil {
		return nil, admin.GenerateUUIDError{Err: err}
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
