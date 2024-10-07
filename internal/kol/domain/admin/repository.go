package admin

import (
	"context"
	"kolresource/internal/kol/entities"
)

type Repository interface {
	GetAdminByUserName(ctx context.Context, userName string) (*entities.Admin, error)
	CreateAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
}
