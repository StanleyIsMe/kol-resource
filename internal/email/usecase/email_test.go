package usecase

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"
	"kolresource/internal/email/mock/repositorymock"
	"kolresource/internal/kol/mock/usecasemock"
	kolUsecase "kolresource/internal/kol/usecase"

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

func TestEmailUseCaseImpl_CreateEmailSender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          CreateEmailSenderParam
	}{
		{
			name:    "success",
			wantErr: false,
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByEmail(gomock.Any(), "test@example.com").Return(nil, commonErrors.ErrDataNotFound)
				repoMock.EXPECT().CreateEmailSender(gomock.Any(), gomock.Any()).Return(nil)

				return repoMock, kolMock
			},
			args: CreateEmailSenderParam{
				UpdatedAdminID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				UpdatedAdminName: "test-admin",
				Name:             "Test Sender",
				Email:            "test@example.com",
				Key:              "test-key",
				RateLimit:        100,
			},
		},
		{
			name:    "duplicate_email",
			wantErr: true,
			expectedError: commonErrors.DuplicatedResourceError{
				Resource: "email sender",
				Name:     "test@example.com",
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByEmail(gomock.Any(), "test@example.com").Return(&entities.EmailSender{}, nil)

				return repoMock, kolMock
			},
			args: CreateEmailSenderParam{
				Email: "test@example.com",
			},
		},
		{
			name:          "GetEmailSenderByEmail_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetEmailSenderByEmail error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: CreateEmailSenderParam{
				Email: "test@example.com",
			},
		},
		{
			name:          "CreateEmailSender_repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.CreateEmailSender error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByEmail(gomock.Any(), gomock.Any()).Return(nil, commonErrors.ErrDataNotFound)
				repoMock.EXPECT().CreateEmailSender(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

				return repoMock, kolMock
			},
			args: CreateEmailSenderParam{
				Email: "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			err := uc.CreateEmailSender(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.CreateEmailSender() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.CreateEmailSender() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_ListEmailSenders(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          []EmailSender
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
	}{
		{
			name:    "success",
			wantErr: false,
			want: []EmailSender{
				{
					ID:         uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:       "Test Sender",
					Email:      "test@example.com",
					RateLimit:  100,
					LastSendAt: fixedTime,
				},
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().AllEmailSenders(gomock.Any()).Return([]*entities.EmailSender{
					{
						ID:         uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name:       "Test Sender",
						Email:      "test@example.com",
						RateLimit:  100,
						LastSendAt: fixedTime,
					},
				}, nil)

				return repoMock, kolMock
			},
		},
		{
			name:    "empty_result",
			wantErr: false,
			want:    []EmailSender{},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().AllEmailSenders(gomock.Any()).Return([]*entities.EmailSender{}, nil)

				return repoMock, kolMock
			},
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListEmailSenders error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().AllEmailSenders(gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			got, err := uc.ListEmailSenders(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.ListEmailSenders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailUseCaseImpl.ListEmailSenders() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.ListEmailSenders() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_UpdateEmailSender(t *testing.T) {
	t.Parallel()

	testName := "Updated Sender"

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          UpdateEmailSenderParam
	}{
		{
			name:    "success",
			wantErr: false,
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().UpdateEmailSender(gomock.Any(), domain.UpdateEmailSenderParam{
					ID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:           &testName,
					UpdatedAdminID: uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				}).Return(nil)

				return repoMock, kolMock
			},
			args: UpdateEmailSenderParam{
				ID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Name:           &testName,
				UpdatedAdminID: uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.UpdateEmailSender error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().UpdateEmailSender(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

				return repoMock, kolMock
			},
			args: UpdateEmailSenderParam{
				ID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			err := uc.UpdateEmailSender(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.UpdateEmailSender() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.UpdateEmailSender() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_GetEmailSender(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          *EmailSender
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          uuid.UUID
	}{
		{
			name:    "success",
			wantErr: false,
			want: &EmailSender{
				ID:         uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Name:       "Test Sender",
				Email:      "test@example.com",
				RateLimit:  100,
				LastSendAt: fixedTime,
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")).Return(&entities.EmailSender{
					ID:         uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:       "Test Sender",
					Email:      "test@example.com",
					RateLimit:  100,
					LastSendAt: fixedTime,
				}, nil)

				return repoMock, kolMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetEmailSenderByID error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			got, err := uc.GetEmailSender(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.GetEmailSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailUseCaseImpl.GetEmailSender() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.GetEmailSender() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_SendEmail(t *testing.T) {
	t.Parallel()

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	productID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")
	adminID := uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6")
	kolID := uuid.MustParse("3193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          SendEmailParam
	}{
		{
			name:    "success",
			wantErr: false,
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(&entities.EmailSender{
					ID:    senderID,
					Name:  "Test Sender",
					Email: "sender@example.com",
				}, nil)

				kolMock.EXPECT().GetProductByID(gomock.Any(), productID).Return(&kolUsecase.Product{
					ID:   productID,
					Name: "Test Product",
				}, nil)

				kolMock.EXPECT().ListKolEmailsByIDs(gomock.Any(), []uuid.UUID{kolID}).Return([]*kolUsecase.KolEmail{
					{
						ID:    kolID,
						Name:  "Test Kol",
						Email: "kol@example.com",
					},
				}, nil)

				repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					},
				)

				repoMock.EXPECT().CreateEmailJob(gomock.Any(), gomock.Any()).Return(&entities.EmailJob{
					ID: 1,
				}, nil)

				repoMock.EXPECT().BatchCreateEmailLogs(gomock.Any(), gomock.Any()).Return(nil)

				return repoMock, kolMock
			},
			args: SendEmailParam{
				Subject:          "Test Subject",
				EmailContent:     "Test Content",
				KolIDs:           []uuid.UUID{kolID},
				ProductID:        productID,
				UpdatedAdminID:   adminID,
				UpdatedAdminName: "test-admin",
				SenderID:         senderID,
			},
		},
		{
			name:    "sender_not_found",
			wantErr: true,
			expectedError: commonErrors.NotFoundError{
				Resource: "email sender",
				ID:       senderID,
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(nil, commonErrors.ErrDataNotFound)

				return repoMock, kolMock
			},
			args: SendEmailParam{
				SenderID: senderID,
			},
		},
		{
			name:          "GetEmailSenderByID_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetEmailSenderByID error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: SendEmailParam{
				SenderID: senderID,
			},
		},
		{
			name:    "product_not_found",
			wantErr: true,
			expectedError: commonErrors.NotFoundError{
				Resource: "product",
				ID:       productID,
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(&entities.EmailSender{
					ID: senderID,
				}, nil)

				kolMock.EXPECT().GetProductByID(gomock.Any(), productID).Return(nil, commonErrors.ErrDataNotFound)

				return repoMock, kolMock
			},
			args: SendEmailParam{
				SenderID:  senderID,
				ProductID: productID,
			},
		},
		{
			name:          "ListKolEmailsByIDs_error",
			wantErr:       true,
			expectedError: fmt.Errorf("kolUsecase.ListKolEmailsByIDs error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(&entities.EmailSender{
					ID: senderID,
				}, nil)

				kolMock.EXPECT().GetProductByID(gomock.Any(), productID).Return(&kolUsecase.Product{
					ID:   productID,
					Name: "Test Product",
				}, nil)

				kolMock.EXPECT().ListKolEmailsByIDs(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: SendEmailParam{
				SenderID:  senderID,
				ProductID: productID,
				KolIDs:    []uuid.UUID{kolID},
			},
		},
		{
			name:    "empty_kols",
			wantErr: true,
			expectedError: commonErrors.NotFoundError{
				Resource: "kols",
				ID:       []uuid.UUID{kolID},
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(&entities.EmailSender{
					ID: senderID,
				}, nil)

				kolMock.EXPECT().GetProductByID(gomock.Any(), productID).Return(&kolUsecase.Product{
					ID:   productID,
					Name: "Test Product",
				}, nil)

				kolMock.EXPECT().ListKolEmailsByIDs(gomock.Any(), []uuid.UUID{kolID}).Return([]*kolUsecase.KolEmail{}, nil)

				return repoMock, kolMock
			},
			args: SendEmailParam{
				SenderID:  senderID,
				ProductID: productID,
				KolIDs:    []uuid.UUID{kolID},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			err := uc.SendEmail(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.SendEmail() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_ListEmailJobs(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          *ListEmailJobsResponse
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          ListEmailJobsParam
	}{
		{
			name:    "success",
			wantErr: false,
			want: &ListEmailJobsResponse{
				EmailJobs: []EmailJob{
					{
						ID:                   1,
						ExpectedReciverCount: 10,
						SuccessCount:         5,
						SenderID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						SenderName:           "Test Sender",
						SenderEmail:          "sender@example.com",
						AdminID:              uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						AdminName:            "Test Admin",
						ProductID:            uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						ProductName:          "Test Product",
						Memo:                 "test memo",
						Status:               email.JobStatusPending,
						CreatedAt:            fixedTime,
						UpdatedAt:            fixedTime,
						LastExecuteAt:        fixedTime,
					},
				},
				Total: 1,
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().ListEmailJobs(gomock.Any(), &domain.ListEmailJobsParams{
					Page: 1,
					Size: 10,
				}).Return([]*entities.EmailJob{
					{
						ID:                   1,
						ExpectedReciverCount: 10,
						SuccessCount:         5,
						SenderID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						SenderName:           "Test Sender",
						SenderEmail:          "sender@example.com",
						AdminID:              uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						AdminName:            "Test Admin",
						ProductID:            uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						ProductName:          "Test Product",
						Memo:                 "test memo",
						Status:               email.JobStatusPending,
						CreatedAt:            fixedTime,
						UpdatedAt:            fixedTime,
						LastExecuteAt:        fixedTime,
					},
				}, int64(1), nil)

				return repoMock, kolMock
			},
			args: ListEmailJobsParam{
				Page:     1,
				PageSize: 10,
			},
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListEmailJobs error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().ListEmailJobs(gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("database error"))

				return repoMock, kolMock
			},
			args: ListEmailJobsParam{
				Page:     1,
				PageSize: 10,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			got, err := uc.ListEmailJobs(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.ListEmailJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailUseCaseImpl.ListEmailJobs() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.ListEmailJobs() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_GetEmailJob(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          *EmailJob
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          int64
	}{
		{
			name:    "success",
			wantErr: false,
			want: &EmailJob{
				ID:                   1,
				ExpectedReciverCount: 10,
				SuccessCount:         5,
				SenderID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				SenderName:           "Test Sender",
				SenderEmail:          "sender@example.com",
				AdminID:              uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				AdminName:            "Test Admin",
				ProductID:            uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				ProductName:          "Test Product",
				Memo:                 "test memo",
				Status:               email.JobStatusPending,
				CreatedAt:            fixedTime,
				UpdatedAt:            fixedTime,
				LastExecuteAt:        fixedTime,
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(&entities.EmailJob{
					ID:                   1,
					ExpectedReciverCount: 10,
					SuccessCount:         5,
					SenderID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					SenderName:           "Test Sender",
					SenderEmail:          "sender@example.com",
					AdminID:              uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					AdminName:            "Test Admin",
					ProductID:            uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					ProductName:          "Test Product",
					Memo:                 "test memo",
					Status:               email.JobStatusPending,
					CreatedAt:            fixedTime,
					UpdatedAt:            fixedTime,
					LastExecuteAt:        fixedTime,
				}, nil)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:    "not_found",
			wantErr: true,
			expectedError: commonErrors.NotFoundError{
				Resource: "email_job",
				ID:       int64(1),
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(nil, commonErrors.ErrDataNotFound)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetEmailJobByID error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			got, err := uc.GetEmailJob(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.GetEmailJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailUseCaseImpl.GetEmailJob() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.GetEmailJob() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_ListEmailLogs(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          []EmailLog
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          ListEmailLogsParam
	}{
		{
			name:    "success",
			wantErr: false,
			want: []EmailLog{
				{
					ID:       1,
					Email:    "kol@example.com",
					Reply:    false,
					KolID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					KolName:  "Test Kol",
					Status:   email.LogStatusPending,
					Memo:     "test memo",
					SendedAt: &fixedTime,
				},
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().ListEmailLogs(gomock.Any(), &domain.ListEmailLogsParams{
					JobID: 1,
				}).Return([]*entities.EmailLog{
					{
						ID:       1,
						Email:    "kol@example.com",
						Reply:    false,
						KolID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						KolName:  "Test Kol",
						Status:   email.LogStatusPending,
						Memo:     "test memo",
						SendedAt: &fixedTime,
					},
				}, nil)

				return repoMock, kolMock
			},
			args: ListEmailLogsParam{
				JobID: 1,
			},
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListEmailLogs error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().ListEmailLogs(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: ListEmailLogsParam{
				JobID: 1,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			got, err := uc.ListEmailLogs(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.ListEmailLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailUseCaseImpl.ListEmailLogs() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.ListEmailLogs() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_CancelEmailJob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          int64
	}{
		{
			name:    "success",
			wantErr: false,
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(&entities.EmailJob{
					ID:     1,
					Status: email.JobStatusPending,
				}, nil)

				repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					},
				)

				repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), int64(1), email.JobStatusCanceled).Return(nil)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:    "not_found",
			wantErr: true,
			expectedError: commonErrors.NotFoundError{
				Resource: "email_job",
				ID:       int64(1),
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(nil, commonErrors.ErrDataNotFound)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:          "invalid_status",
			wantErr:       true,
			expectedError: fmt.Errorf("email job status is not cancelable"),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(&entities.EmailJob{
					ID:     1,
					Status: email.JobStatusSuccess,
				}, nil)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetEmailJobByID error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			err := uc.CancelEmailJob(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.CancelEmailJob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.CancelEmailJob() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestEmailUseCaseImpl_StartEmailJob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getMocks      func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase)
		args          int64
	}{
		{
			name:    "success",
			wantErr: false,
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(&entities.EmailJob{
					ID:     1,
					Status: email.JobStatusCanceled,
				}, nil)

				repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					},
				)

				repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), int64(1), email.JobStatusPending).Return(nil)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:    "not_found",
			wantErr: true,
			expectedError: commonErrors.NotFoundError{
				Resource: "email_job",
				ID:       int64(1),
			},
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(nil, commonErrors.ErrDataNotFound)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:          "invalid_status",
			wantErr:       true,
			expectedError: fmt.Errorf("email job status is not startable"),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), int64(1)).Return(&entities.EmailJob{
					ID:     1,
					Status: email.JobStatusPending,
				}, nil)

				return repoMock, kolMock
			},
			args: 1,
		},
		{
			name:          "repo_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetEmailJobByID error: %w", errors.New("database error")),
			getMocks: func(ctrl *gomock.Controller) (domain.Repository, kolUsecase.KolUseCase) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				kolMock := usecasemock.NewMockKolUseCase(ctrl)

				repoMock.EXPECT().GetEmailJobByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, kolMock
			},
			args: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, kolMock := tt.getMocks(ctrl)
			uc := NewEmailUseCaseImpl(repoMock, kolMock)

			err := uc.StartEmailJob(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailUseCaseImpl.StartEmailJob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("EmailUseCaseImpl.StartEmailJob() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
