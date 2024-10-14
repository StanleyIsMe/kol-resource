package usecase

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AdminUseCase interface {
	Register(ctx context.Context, param RegisterParams) error
	Login(ctx context.Context, userName, password string) (*LoginResponse, error)
	LoginTokenParser(ctx context.Context, tokenString string) (*JWTAdminClaims, error)
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

type JWTAdminClaims struct {
	AdminID   uuid.UUID `json:"admin_id"`
	AdminName string    `json:"admin_name"`
	jwt.RegisteredClaims
}
