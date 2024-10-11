package usecase

import (
	"context"
	"kolresource/internal/kol"

	"github.com/google/uuid"
)

type KolUseCase interface {
	CreateKol(ctx context.Context, param CreateKolParam) error
	GetKolByID(ctx context.Context, kolID uuid.UUID) (*Kol, error)
	UpdateKol(ctx context.Context, param UpdateKolParam) error
	ListKols(ctx context.Context, param ListKolsParam) ([]*Kol, int, error)

	CreateTag(ctx context.Context, param CreateTagParam) error
	ListTagsByName(ctx context.Context, name string) ([]*Tag, error)

	CreateProduct(ctx context.Context, param CreateProductParam) error
	ListProductsByName(ctx context.Context, name string) ([]*Product, error)

	SendEmail(ctx context.Context, param SendEmailParam) error
}

type Kol struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
	Sex         kol.Sex   `json:"sex"`
	Tags        []Tag     `json:"tags"`
}

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type CreateKolParam struct {
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	Description    string      `json:"description"`
	Sex            kol.Sex     `json:"sex"`
	Tags           []uuid.UUID `json:"tags"`
	UpdatedAdminID uuid.UUID   `json:"updated_admin_id"`
}

type UpdateKolParam struct {
	KolID          uuid.UUID   `json:"kol_id"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	Description    string      `json:"description"`
	Sex            kol.Sex     `json:"sex"`
	Tags           []uuid.UUID `json:"tags"`
	UpdatedAdminID uuid.UUID   `json:"updated_admin_id"`
}

type ListKolsParam struct {
	Email    *string
	Name     *string
	Tag      *string
	Sex      *kol.Sex
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type CreateTagParam struct {
	Name           string    `json:"name"`
	UpdatedAdminID uuid.UUID `json:"updated_admin_id"`
}

type CreateProductParam struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	UpdatedAdminID uuid.UUID `json:"updated_admin_id"`
}

type SendEmailParam struct {
	EmailContent string      `json:"email_content"`
	KolIDs       []uuid.UUID `json:"kol_ids"`
}
