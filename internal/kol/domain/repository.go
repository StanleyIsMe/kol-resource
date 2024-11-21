package domain

import (
	"context"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain/entities"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=../mock/repositorymock/repository_mock.go -package=repositorymock
//nolint:interfacebloat
type Repository interface {
	GetKolByID(ctx context.Context, id uuid.UUID) (*entities.Kol, error)
	GetKolByEmail(ctx context.Context, email string) (*entities.Kol, error)
	CreateKol(ctx context.Context, param CreateKolParams) (*entities.Kol, error)
	UpdateKol(ctx context.Context, param UpdateKolParams) (*entities.Kol, error)
	DeleteKolByID(ctx context.Context, id uuid.UUID) error
	GetKolWithTagsByID(ctx context.Context, id uuid.UUID) (*Kol, error)
	ListKolWithTagsByFilters(ctx context.Context, param ListKolWithTagsByFiltersParams) ([]*Kol, int, error)
	ListKolsByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.Kol, error)

	CreateTag(ctx context.Context, param CreateTagParams) (*entities.Tag, error)
	ListTagsByName(ctx context.Context, name string) ([]*entities.Tag, error)
	DeleteTagByID(ctx context.Context, id uuid.UUID) error
	GetTagByName(ctx context.Context, name string) (*entities.Tag, error)

	CreateProduct(ctx context.Context, param CreateProductParams) (*entities.Product, error)
	ListProductsByName(ctx context.Context, name string) ([]*entities.Product, error)
	GetProductByName(ctx context.Context, name string) (*entities.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	DeleteProductByID(ctx context.Context, id uuid.UUID) error

	CreateSendEmailLog(ctx context.Context, sendEmailLog *entities.SendEmailLog) (*entities.SendEmailLog, error)
	ListSendEmailLogsByFilter(ctx context.Context, param ListSendEmailLogsByFilterParams) ([]*entities.SendEmailLog, int, error)
}

type EmailRepository interface {
	SendEmail(ctx context.Context, param SendEmailParams) error
}

type CreateKolParams struct {
	Name           string
	Email          string
	SocialMedia    string
	Description    string
	Sex            kol.Sex
	Enable         bool
	UpdatedAdminID uuid.UUID
	Tags           []uuid.UUID
}

type Tag struct {
	ID   uuid.UUID
	Name string
}

type UpdateKolParams struct {
	ID             uuid.UUID
	Name           string
	Email          string
	SocialMedia    string
	Description    string
	Sex            kol.Sex
	Enable         bool
	UpdatedAdminID uuid.UUID
	Tags           []uuid.UUID
}

type ListKolWithTagsByFiltersParams struct {
	Email    *string
	Name     *string
	Tag      *string
	TagIDs   []uuid.UUID
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

type CreateSendEmailLogParams struct {
	KolID       uuid.UUID
	KolName     string
	Email       string
	AdminID     uuid.UUID
	AdminName   string
	ProductID   uuid.UUID
	ProductName string
}

type SendEmailParams struct {
	AdminEmail string
	AdminPass  string
	Subject    string
	Body       string
	ToEmails   []ToEmail
}

type ToEmail struct {
	Email string
	Name  string
}
