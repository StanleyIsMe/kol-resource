package domain

import (
	"context"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain/entities"

	"github.com/google/uuid"
)

type Repository interface {
	GetKolByID(ctx context.Context, id uuid.UUID) (*entities.Kol, error)
	CreateKol(ctx context.Context, param CreateKolParams) (*entities.Kol, error)
	UpdateKol(ctx context.Context, param UpdateKolParams) (*entities.Kol, error)
	DeleteKolByID(ctx context.Context, id uuid.UUID) error
	GetKolWithTagByID(ctx context.Context, id uuid.UUID) (*Kol, error)
	ListKolWithTagsByFilters(ctx context.Context, param ListKolWithTagsByFiltersParams) ([]*Kol, int, error)

	CreateTag(ctx context.Context, param CreateTagParams) (*entities.Tag, error)
	ListTagsByName(ctx context.Context, name string) ([]*entities.Tag, error)
	DeleteTagByID(ctx context.Context, id uuid.UUID) error

	CreateProduct(ctx context.Context, param CreateProductParams) (*entities.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	DeleteProductByID(ctx context.Context, id uuid.UUID) error

	CreateSendEmailLog(ctx context.Context, sendEmailLog *entities.SendEmailLog) (*entities.SendEmailLog, error)
	ListSendEmailLogsByFilter(ctx context.Context, param ListSendEmailLogsByFilterParams) ([]*entities.SendEmailLog, int, error)
}

type CreateKolParams struct {
	Name           string
	Email          string
	Description    string
	Sex            kol.Sex
	Enable         bool
	UpdatedAdminID uuid.UUID
}

type UpdateKolParams struct {
	ID             uuid.UUID
	Name           string
	Email          string
	Description    string
	Sex            kol.Sex
	Enable         bool
	UpdatedAdminID uuid.UUID
}

type ListKolWithTagsByFiltersParams struct {
	Email    *string
	Name     *string
	Tag      *string
	Sex      *kol.Sex
	Page     int
	PageSize int
}

type CreateTagParams struct {
	Name           string
	UpdatedAdminID uuid.UUID
}

type CreateProductParams struct {
	Name           string
	Description    string
	UpdatedAdminID uuid.UUID
}

type ListSendEmailLogsByFilterParams struct {
	Email       *string
	ProductName *string
	AdminName   *string
	KolName     *string
	Page        int
	PageSize    int
}
