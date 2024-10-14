package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"kolresource/internal/admin/domain"
	"kolresource/internal/admin/domain/entities"

	apiCfg "kolresource/internal/api/config"
	"kolresource/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AdminUseCaseImpl struct {
	adminRepo domain.Repository
	cfg       *config.Config[apiCfg.Config]
}

func NewAdminUseCaseImpl(adminRepo domain.Repository, cfg *config.Config[apiCfg.Config]) *AdminUseCaseImpl {
	return &AdminUseCaseImpl{adminRepo: adminRepo, cfg: cfg}
}

// Register is responsible for registering a new admin user.
func (a *AdminUseCaseImpl) Register(ctx context.Context, param RegisterParams) error {
	existAdmin, err := a.adminRepo.GetAdminByUserName(ctx, param.UserName)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		return InternalServerError{err: fmt.Errorf("adminRepo.GetAdminByUserName error: %w", err)}
	}

	if existAdmin != nil {
		return DumplicatedUsernameError{username: param.UserName}
	}

	argon2IDHash := NewArgon2idHash(1, 32, 64*1024, 1, 128)

	hashSalt, err := argon2IDHash.GenerateHash([]byte(param.Password), nil)
	if err != nil {
		return InternalServerError{err: fmt.Errorf("argon2IDHash.GenerateHash error: %w", err)}
	}

	adminEntity := &entities.Admin{
		Username: param.UserName,
		Name:     param.Name,
		Salt:     base64.StdEncoding.EncodeToString(hashSalt.Salt),
		Password: base64.StdEncoding.EncodeToString(hashSalt.Hash),
	}

	if _, err := a.adminRepo.CreateAdmin(ctx, adminEntity); err != nil {
		return InternalServerError{err: fmt.Errorf("adminRepo.CreateAdmin error: %w", err)}
	}

	return nil
}

// Login is responsible for logging in an admin user.
func (a *AdminUseCaseImpl) Login(ctx context.Context, userName string, password string) (*LoginResponse, error) {
	adminEntity, err := a.adminRepo.GetAdminByUserName(ctx, userName)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, UnauthorizedError{err: errors.New("username not found")}
		}

		return nil, InternalServerError{err: fmt.Errorf("adminRepo.GetAdminByUserName error: %w", err)}
	}

	argon2IDHash := NewArgon2idHash(1, 32, 64*1024, 1, 128)

	hash, err := base64.StdEncoding.DecodeString(adminEntity.Password)
	if err != nil {
		return nil, InternalServerError{err: fmt.Errorf("base64.StdEncoding.DecodeString error: %w", err)}
	}

	salt, err := base64.StdEncoding.DecodeString(adminEntity.Salt)
	if err != nil {
		return nil, InternalServerError{err: fmt.Errorf("base64.StdEncoding.DecodeString error: %w", err)}
	}

	if err := argon2IDHash.Compare(hash, salt, []byte(password)); err != nil {
		return nil, UnauthorizedError{err: fmt.Errorf("argon2IDHash.Compare error: %w", err)}
	}

	token, err := a.generateJWT(adminEntity.ID, adminEntity.Name)
	if err != nil {
		return nil, InternalServerError{err: fmt.Errorf("generateJWT error: %w", err)}
	}

	return &LoginResponse{
		Token:     token,
		AdminName: adminEntity.Name,
	}, nil
}

const (
	signKey = "kolresourceKey"
)

func (a *AdminUseCaseImpl) generateJWT(adminID uuid.UUID, adminName string) (string, error) {
	claims := JWTAdminClaims{
		adminID,
		adminName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.cfg.CustomConfig.Auth.JWTExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "kolresource",
			Subject:   "admin",
			ID:        "1",
			Audience:  []string{"stanley"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString([]byte(a.cfg.CustomConfig.Auth.JWTKey))
	if err != nil {
		return "", fmt.Errorf("jwt generate error: %w", err)
	}

	return jwtToken, nil
}

func (a *AdminUseCaseImpl) LoginTokenParser(ctx context.Context, tokenString string) (*JWTAdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTAdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.CustomConfig.Auth.JWTKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt parse error: %w", err)
	}

	claims, ok := token.Claims.(*JWTAdminClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
