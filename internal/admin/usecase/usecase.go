package usecase

import "context"

type AdminUseCase interface {
	Register(ctx context.Context, param RegisterParams) error
	Login(ctx context.Context, userName, password string) (*LoginResponse, error)
	LoginTokenParser(ctx context.Context, tokenString string) error
}

type RegisterParams struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	AdminName string `json:"admin_name"`
}
