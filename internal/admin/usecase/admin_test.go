package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"os"
	"testing"
	"time"

	"kolresource/internal/admin/domain"
	"kolresource/internal/admin/domain/entities"
	"kolresource/internal/admin/mock/repositorymock"
	apiCfg "kolresource/internal/api/config"
	"kolresource/pkg/config"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func newTestConfig() *config.Config[apiCfg.Config] {
	return &config.Config[apiCfg.Config]{
		CustomConfig: apiCfg.Config{
			Auth: apiCfg.Auth{
				JWTKey: "test-secret-key-for-jwt-signing",
				JWTExp: 1 * time.Hour,
			},
		},
	}
}

func TestAdminUseCaseImpl_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          RegisterParams
	}{
		{
			name:    "success",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(&entities.Admin{}, nil)

				return repoMock
			},
			args: RegisterParams{
				Name:     "test",
				UserName: "testuser",
				Password: "testpassword",
			},
		},
		{
			name:          "duplicate username",
			wantErr:       true,
			expectedError: DumplicatedUsernameError{username: "testuser"},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(&entities.Admin{
					ID:       uuid.New(),
					Username: "testuser",
				}, nil)

				return repoMock
			},
			args: RegisterParams{
				Name:     "test",
				UserName: "testuser",
				Password: "testpassword",
			},
		},
		{
			name:          "GetAdminByUserName error",
			wantErr:       true,
			expectedError: InternalServerError{err: errors.New("adminRepo.GetAdminByUserName error: db connection failed")},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(nil, errors.New("db connection failed"))

				return repoMock
			},
			args: RegisterParams{
				Name:     "test",
				UserName: "testuser",
				Password: "testpassword",
			},
		},
		{
			name:          "CreateAdmin error",
			wantErr:       true,
			expectedError: InternalServerError{err: errors.New("adminRepo.CreateAdmin error: insert failed")},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(nil, errors.New("insert failed"))

				return repoMock
			},
			args: RegisterParams{
				Name:     "test",
				UserName: "testuser",
				Password: "testpassword",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := tt.getRepoMock(ctrl)
			cfg := newTestConfig()

			uc := NewAdminUseCaseImpl(repoMock, cfg)
			err := uc.Register(context.Background(), tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.expectedError != nil {
					gotType := errors.Unwrap(err)
					expectedType := errors.Unwrap(tt.expectedError)

					if gotType != nil && expectedType != nil {
						return
					}

					switch tt.expectedError.(type) {
					case InternalServerError:
						var target InternalServerError
						if !errors.As(err, &target) {
							t.Errorf("Register() error type = %T, want InternalServerError", err)
						}
					case DumplicatedUsernameError:
						var target DumplicatedUsernameError
						if !errors.As(err, &target) {
							t.Errorf("Register() error type = %T, want DumplicatedUsernameError", err)
						}
					}
				}
			}
		})
	}
}

func TestAdminUseCaseImpl_Login(t *testing.T) {
	t.Parallel()

	argon2Hash := NewArgon2idHash(1, 32, 64*1024, 1, 128)
	hashSalt, _ := argon2Hash.GenerateHash([]byte("testpassword"), nil)

	admin := &entities.Admin{
		ID:       uuid.New(),
		Name:     "test",
		Username: "testuser",
		Password: base64.StdEncoding.EncodeToString(hashSalt.Hash),
		Salt:     base64.StdEncoding.EncodeToString(hashSalt.Salt),
	}

	tests := []struct {
		name          string
		userName      string
		password      string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
	}{
		{
			name:     "success",
			userName: "testuser",
			password: "testpassword",
			wantErr:  false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(admin, nil)

				return repoMock
			},
		},
		{
			name:          "user not found",
			userName:      "unknownuser",
			password:      "testpassword",
			wantErr:       true,
			expectedError: UnauthorizedError{err: errors.New("username not found")},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "unknownuser").Return(nil, domain.ErrDataNotFound)

				return repoMock
			},
		},
		{
			name:          "wrong password",
			userName:      "testuser",
			password:      "wrongpassword",
			wantErr:       true,
			expectedError: UnauthorizedError{err: errors.New("argon2IDHash.Compare error: hash doesn't match")},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(admin, nil)

				return repoMock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := tt.getRepoMock(ctrl)
			cfg := newTestConfig()

			uc := NewAdminUseCaseImpl(repoMock, cfg)
			resp, err := uc.Login(context.Background(), tt.userName, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("Login() response should not be nil on success")
					return
				}

				if resp.Token == "" {
					t.Error("Login() token should not be empty")
				}

				if resp.AdminName != admin.Name {
					t.Errorf("Login() adminName = %v, want %v", resp.AdminName, admin.Name)
				}
			}

			if tt.wantErr && tt.expectedError != nil {
				switch tt.expectedError.(type) {
				case UnauthorizedError:
					var target UnauthorizedError
					if !errors.As(err, &target) {
						t.Errorf("Login() error type = %T, want UnauthorizedError", err)
					}
				case InternalServerError:
					var target InternalServerError
					if !errors.As(err, &target) {
						t.Errorf("Login() error type = %T, want InternalServerError", err)
					}
				}
			}
		})
	}
}

func TestAdminUseCaseImpl_Login_InvalidBase64Password(t *testing.T) {
	t.Parallel()

	adminWithBadPassword := &entities.Admin{
		ID:       uuid.New(),
		Name:     "test",
		Username: "testuser",
		Password: "!!!not-valid-base64!!!",
		Salt:     base64.StdEncoding.EncodeToString([]byte("valid-salt")),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repositorymock.NewMockRepository(ctrl)
	repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(adminWithBadPassword, nil)

	cfg := newTestConfig()
	uc := NewAdminUseCaseImpl(repoMock, cfg)

	_, err := uc.Login(context.Background(), "testuser", "testpassword")
	if err == nil {
		t.Fatal("Login() expected error for invalid base64 password, got nil")
	}

	var target InternalServerError
	if !errors.As(err, &target) {
		t.Errorf("Login() error type = %T, want InternalServerError", err)
	}
}

func TestAdminUseCaseImpl_Login_InvalidBase64Salt(t *testing.T) {
	t.Parallel()

	adminWithBadSalt := &entities.Admin{
		ID:       uuid.New(),
		Name:     "test",
		Username: "testuser",
		Password: base64.StdEncoding.EncodeToString([]byte("valid-hash")),
		Salt:     "!!!not-valid-base64!!!",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repositorymock.NewMockRepository(ctrl)
	repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(adminWithBadSalt, nil)

	cfg := newTestConfig()
	uc := NewAdminUseCaseImpl(repoMock, cfg)

	_, err := uc.Login(context.Background(), "testuser", "testpassword")
	if err == nil {
		t.Fatal("Login() expected error for invalid base64 salt, got nil")
	}

	var target InternalServerError
	if !errors.As(err, &target) {
		t.Errorf("Login() error type = %T, want InternalServerError", err)
	}
}

func TestAdminUseCaseImpl_LoginTokenParser(t *testing.T) {
	t.Parallel()

	argon2Hash := NewArgon2idHash(1, 32, 64*1024, 1, 128)
	hashSalt, _ := argon2Hash.GenerateHash([]byte("testpassword"), nil)

	adminID := uuid.New()
	admin := &entities.Admin{
		ID:       adminID,
		Name:     "test",
		Username: "testuser",
		Password: base64.StdEncoding.EncodeToString(hashSalt.Hash),
		Salt:     base64.StdEncoding.EncodeToString(hashSalt.Salt),
	}

	cfg := newTestConfig()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repositorymock.NewMockRepository(ctrl)
	repoMock.EXPECT().GetAdminByUserName(gomock.Any(), "testuser").Return(admin, nil)

	uc := NewAdminUseCaseImpl(repoMock, cfg)
	loginResp, err := uc.Login(context.Background(), "testuser", "testpassword")

	if err != nil {
		t.Fatalf("Login() failed to generate token for LoginTokenParser test: %v", err)
	}

	validToken := loginResp.Token

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
		checkClaims func(t *testing.T, claims *JWTAdminClaims)
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			wantErr:     false,
			checkClaims: func(t *testing.T, claims *JWTAdminClaims) {
				t.Helper()

				if claims.AdminID != adminID {
					t.Errorf("LoginTokenParser() AdminID = %v, want %v", claims.AdminID, adminID)
				}

				if claims.AdminName != "test" {
					t.Errorf("LoginTokenParser() AdminName = %v, want %v", claims.AdminName, "test")
				}
			},
		},
		{
			name:        "invalid token string",
			tokenString: "invalid.token.string",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			claims, err := uc.LoginTokenParser(context.Background(), tt.tokenString)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoginTokenParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkClaims != nil {
				tt.checkClaims(t, claims)
			}
		})
	}
}
