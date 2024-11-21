package domain

import (
	"context"
	"kolresource/internal/admin/domain/entities"
)

//go:generate mockgen -source=./repository.go -destination=../mock/repositorymock/repository_mock.go -package=repositorymock
type Repository interface {
	GetAdminByUserName(ctx context.Context, userName string) (*entities.Admin, error)
	CreateAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
}
