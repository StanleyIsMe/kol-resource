package usecase

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain"
	"kolresource/internal/kol/domain/entities"
	"kolresource/internal/kol/mock/repositorymock"
	"os"
	"reflect"
	"testing"

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

func TestKolUseCaseImpl_CreateKol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expectedError error
		wantErr       bool
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          CreateKolParam
	}{
		{
			name:    "success",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "test@example.com").Return(nil, domain.ErrDataNotFound)

				createKolParams := domain.CreateKolParams{
					Name:           "test",
					Email:          "test@example.com",
					Description:    "test",
					SocialMedia:    "test",
					Sex:            kol.SexMale,
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Tags: []uuid.UUID{
						uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					},
					Enable: true,
				}

				repoMock.EXPECT().CreateKol(gomock.Any(), createKolParams).Return(nil, nil)

				return repoMock
			},
			args: CreateKolParam{
				Name:           "test",
				Email:          "test@example.com",
				Description:    "test",
				SocialMedia:    "test",
				Sex:            kol.SexMale,
				UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Tags: []uuid.UUID{
					uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				},
			},
		},
		{
			name:    "duplicated_email",
			wantErr: true,
			expectedError: DuplicatedResourceError{
				resource: "kol",
				name:     "test@example.com",
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "test@example.com").Return(&entities.Kol{}, nil)

				return repoMock
			},
			args: CreateKolParam{
				Email: "test@example.com",
			},
		},
		{
			name:          "GetKolByEmail_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetKolByEmail error: %w", errors.New("test")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "test@example.com").Return(nil, errors.New("test"))

				return repoMock
			},
			args: CreateKolParam{
				Email: "test@example.com",
			},
		},
		{
			name:          "CreateKol_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.CreateKol error: %w", errors.New("test")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateKol(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))

				return repoMock
			},
			args: CreateKolParam{
				Email: "test@example.com",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			err := uc.CreateKol(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.CreateKol() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.CreateKol() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_GetKolByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          *Kol
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          uuid.UUID
	}{
		{
			name:    "success",
			wantErr: false,
			want: &Kol{
				ID:          uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Name:        "Test Name",
				Email:       "test@example.com",
				Description: "Test Description",
				SocialMedia: "Test Social",
				Sex:         kol.SexMale,
				Tags: []Tag{
					{
						ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name: "Tag1",
					},
				},
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				kolAggregate := domain.NewKol(&entities.Kol{
					ID:          uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:        "Test Name",
					Email:       "test@example.com",
					Description: "Test Description",
					SocialMedia: "Test Social",
					Sex:         kol.SexMale,
				})
				kolAggregate.AppendTag(&entities.Tag{
					ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name: "Tag1",
				})

				repoMock.EXPECT().GetKolWithTagsByID(gomock.Any(), uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")).Return(kolAggregate, nil)

				return repoMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "not_found",
			wantErr:       true,
			expectedError: NotFoundError{resource: "kol", id: "0193487b-f1a2-7a72-8ae4-197b84dc52d6"},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolWithTagsByID(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)

				return repoMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "GetKolWithTagsByID_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetKolWithTagsByID error: %w", errors.New("test")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolWithTagsByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))

				return repoMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			got, err := uc.GetKolByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.GetKolByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KolUseCaseImpl.GetKolByID() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.GetKolByID() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_UpdateKol(t *testing.T) {
	tests := []struct {
		name          string
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          UpdateKolParam
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				updateKolParams := domain.UpdateKolParams{
					ID:             uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:           "Test Name",
					Email:          "test@email.com",
					Description:    "Test Description",
					SocialMedia:    "Test Social",
					Sex:            kol.SexMale,
					Tags:           []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				}

				repoMock.EXPECT().UpdateKol(gomock.Any(), updateKolParams).Return(nil, nil)

				return repoMock
			},
			args: UpdateKolParam{
				KolID:          uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Name:           "Test Name",
				Email:          "test@email.com",
				Description:    "Test Description",
				SocialMedia:    "Test Social",
				Sex:            kol.SexMale,
				Tags:           []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
				UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "not_found_error",
			wantErr:       true,
			expectedError: NotFoundError{resource: "kol", id: "0193487b-f1a2-7a72-8ae4-197b84dc52d6"},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().UpdateKol(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)

				return repoMock
			},
			args: UpdateKolParam{
				KolID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "update_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.UpdateKol error: %w", errors.New("test")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().UpdateKol(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))

				return repoMock
			},
			args: UpdateKolParam{
				KolID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			err := uc.UpdateKol(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.UpdateKol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.UpdateKol() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_ListKols(t *testing.T) {
	tests := []struct {
		name          string
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          ListKolsParam
		want          []*Kol
		wantTotal     int
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				kolID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
				tagID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

				kolAggregate := domain.NewKol(&entities.Kol{
					ID:          kolID,
					Name:        "Test Name",
					Email:       "test@email.com",
					Description: "Test Description",
					SocialMedia: "Test Social",
					Sex:         kol.SexMale,
					Enable:      true,
				})
				kolAggregate.AppendTag(&entities.Tag{
					ID:   tagID,
					Name: "Test Tag",
				})

				expectedParams := domain.ListKolWithTagsByFiltersParams{
					TagIDs:   []uuid.UUID{tagID},
					Page:     1,
					PageSize: 10,
				}

				repoMock.EXPECT().ListKolWithTagsByFilters(gomock.Any(), expectedParams).Return([]*domain.Kol{kolAggregate}, 1, nil)

				return repoMock
			},
			args: ListKolsParam{
				TagIDs:   []uuid.UUID{uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")},
				Page:     1,
				PageSize: 10,
			},
			want: []*Kol{
				{
					ID:          uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:        "Test Name",
					Email:       "test@email.com",
					Description: "Test Description",
					SocialMedia: "Test Social",
					Sex:         kol.SexMale,
					Tags: []Tag{
						{
							ID:   uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"),
							Name: "Test Tag",
						},
					},
				},
			},
			wantTotal: 1,
		},
		{
			name: "empty_result",
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				expectedParams := domain.ListKolWithTagsByFiltersParams{
					Page:     1,
					PageSize: 10,
				}

				repoMock.EXPECT().ListKolWithTagsByFilters(gomock.Any(), expectedParams).Return([]*domain.Kol{}, 0, nil)

				return repoMock
			},
			args: ListKolsParam{
				Page:     1,
				PageSize: 10,
			},
			want:      []*Kol{},
			wantTotal: 0,
		},
		{
			name:          "list_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListKolWithTagsByFilters error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)
				repoMock.EXPECT().ListKolWithTagsByFilters(gomock.Any(), gomock.Any()).Return(nil, 0, errors.New("database error"))
				return repoMock
			},
			args: ListKolsParam{
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

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			got, total, err := uc.ListKols(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.ListKols() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("KolUseCaseImpl.ListKols() got = %v, want %v", got, tt.want)
				}
				if total != tt.wantTotal {
					t.Errorf("KolUseCaseImpl.ListKols() total = %v, want %v", total, tt.wantTotal)
				}
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.ListKols() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_CreateTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          CreateTagParam
	}{
		{
			name:    "success",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetTagByName(gomock.Any(), "test-tag").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateTag(gomock.Any(), domain.CreateTagParams{
					Name:           "test-tag",
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				}).Return(&entities.Tag{}, nil)

				return repoMock
			},
			args: CreateTagParam{
				Name:           "test-tag",
				UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "duplicate_tag",
			wantErr:       true,
			expectedError: DuplicatedResourceError{resource: "tag", name: "test-tag"},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetTagByName(gomock.Any(), "test-tag").Return(&entities.Tag{}, nil)

				return repoMock
			},
			args: CreateTagParam{
				Name: "test-tag",
			},
		},
		{
			name:          "get_tag_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetTagByName error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: CreateTagParam{
				Name: "test-tag",
			},
		},
		{
			name:          "create_tag_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.CreateTag error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: CreateTagParam{
				Name: "test-tag",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			err := uc.CreateTag(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.CreateTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.CreateTag() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_ListTagsByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          []*Tag
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          string
	}{
		{
			name:    "success",
			wantErr: false,
			want: []*Tag{
				{
					ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name: "test-tag",
				},
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().ListTagsByName(gomock.Any(), "test").Return([]*entities.Tag{
					{
						ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name: "test-tag",
					},
				}, nil)

				return repoMock
			},
			args: "test",
		},
		{
			name:    "empty_result",
			wantErr: false,
			want:    []*Tag{},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().ListTagsByName(gomock.Any(), gomock.Any()).Return([]*entities.Tag{}, nil)

				return repoMock
			},
			args: "nonexistent",
		},
		{
			name:          "list_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListTagsByName error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().ListTagsByName(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			got, err := uc.ListTagsByName(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.ListTagsByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KolUseCaseImpl.ListTagsByName() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.ListTagsByName() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_ListProductsByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		want          []*Product
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          string
	}{
		{
			name:    "success",
			wantErr: false,
			want: []*Product{
				{
					ID:          uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name:        "test-product",
					Description: "test description",
				},
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().ListProductsByName(gomock.Any(), "test").Return([]*entities.Product{
					{
						ID:          uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name:        "test-product",
						Description: "test description",
					},
				}, nil)

				return repoMock
			},
			args: "test",
		},
		{
			name:    "empty_result",
			wantErr: false,
			want:    []*Product{},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().ListProductsByName(gomock.Any(), gomock.Any()).Return(nil, nil)

				return repoMock
			},
			args: "nonexistent",
		},
		{
			name:          "list_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListProductsByName error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().ListProductsByName(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			got, err := uc.ListProductsByName(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.ListProductsByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KolUseCaseImpl.ListProductsByName() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.ListProductsByName() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}


func TestKolUseCaseImpl_CreateProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          CreateProductParam
	}{
		{
			name:    "success",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetProductByName(gomock.Any(), "test-product").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateProduct(gomock.Any(), domain.CreateProductParams{
					Name:           "test-product",
					Description:    "test description",
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				}).Return(&entities.Product{}, nil)

				return repoMock
			},
			args: CreateProductParam{
				Name:           "test-product",
				Description:    "test description",
				UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "duplicate_product",
			wantErr:       true,
			expectedError: DuplicatedResourceError{resource: "product", name: "test-product"},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetProductByName(gomock.Any(), "test-product").Return(&entities.Product{}, nil)

				return repoMock
			},
			args: CreateProductParam{
				Name: "test-product",
			},
		},
		{
			name:          "get_product_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetProductByName error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetProductByName(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: CreateProductParam{
				Name: "test-product",
			},
		},
		{
			name:          "create_product_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.CreateProduct error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetProductByName(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateProduct(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: CreateProductParam{
				Name: "test-product",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			err := uc.CreateProduct(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.CreateProduct() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_SendEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository)
		args          SendEmailParam
	}{
		{
			name:    "success",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")).Return(&entities.Product{
					ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name: "test-product",
				}, nil)

				kols := []*entities.Kol{
					{
						ID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name:  "test-kol",
						Email: "test@example.com",
					},
				}

				repoMock.EXPECT().ListKolsByIDs(gomock.Any(), []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")}).Return(kols, nil)

				sendEmailParams := domain.SendEmailParams{
					Subject:    "Test Subject",
					Body:       "Test Content",
					ToEmails:   []domain.ToEmail{{Email: "test@example.com", Name: "test-kol"}},
					Images:     []domain.SendEmailImage{{ContentID: "test-content-id", Data: "test-data", ImageType: "test-image-type"}},
				}

				emailRepoMock.EXPECT().SendEmail(gomock.Any(), sendEmailParams).Return(nil)

				createSendEmailLogParam := &entities.SendEmailLog{
					AdminID:     uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					AdminName:   "test-admin",
					KolID:       uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					KolName:     "test-kol",
					Email:       "test@example.com",
					ProductID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					ProductName: "test-product",
				}
				repoMock.EXPECT().CreateSendEmailLog(gomock.Any(), createSendEmailLogParam).Return(nil, nil)

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				KolIDs:       []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
				ProductID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Subject:      "Test Subject",
				EmailContent: "Test Content",
				UpdatedAdminID:     uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				UpdatedAdminName:   "test-admin",
				Images:           []SendEmailImage{{ContentID: "test-content-id", Data: "test-data", ImageType: "test-image-type"}},
			},
		},
		{
			name:          "GetProductByID_not_found",
			wantErr:       true,
			expectedError: NotFoundError{resource: "product", id: "0193487b-f1a2-7a72-8ae4-197b84dc52d6"},
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				ProductID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "GetProductByID_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetProductByID error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				ProductID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
		},
		{
			name:          "ListKolsByIDs_not_found",
			wantErr:       true,
			expectedError: NotFoundError{resource: "kol", id: []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")}},
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(nil, nil)

				repoMock.EXPECT().ListKolsByIDs(gomock.Any(), gomock.Any()).Return(nil, nil)

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				KolIDs: []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
			},
		},
		{
			name:          "ListKolsByIDs_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.ListKolsByIDs error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(nil, nil)
				repoMock.EXPECT().ListKolsByIDs(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				KolIDs: []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
			},
		},
		{
			name:          "send_email_error",
			wantErr:       true,
			expectedError: fmt.Errorf("emailRepo.SendEmail error: %w", errors.New("email error")),
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")).Return(&entities.Product{
					ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name: "test-product",
				}, nil)

				kols := []*entities.Kol{
					{
						ID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name:  "test-kol",
						Email: "test@example.com",
					},
				}

				repoMock.EXPECT().ListKolsByIDs(gomock.Any(), []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")}).Return(kols, nil)
				emailRepoMock.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(errors.New("email error"))

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				KolIDs:       []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
				ProductID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Subject:      "Test Subject",
				EmailContent: "Test Content",
				UpdatedAdminID:     uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				UpdatedAdminName:   "test-admin",
			},
		},
		{
			name:          "create_send_email_log_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.CreateSendEmailLog error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) (domain.Repository, domain.EmailRepository) {
				repoMock := repositorymock.NewMockRepository(ctrl)
				emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

				repoMock.EXPECT().GetProductByID(gomock.Any(), uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")).Return(&entities.Product{
					ID:   uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Name: "test-product",
				}, nil)

				kols := []*entities.Kol{
					{
						ID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name:  "test-kol",
						Email: "test@example.com",
					},
				}

				repoMock.EXPECT().ListKolsByIDs(gomock.Any(), []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")}).Return(kols, nil)
				emailRepoMock.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil)

				repoMock.EXPECT().CreateSendEmailLog(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock, emailRepoMock
			},
			args: SendEmailParam{
				KolIDs:       []uuid.UUID{uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")},
				ProductID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Subject:      "Test Subject",
				EmailContent: "Test Content",
				UpdatedAdminID:     uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				UpdatedAdminName:   "test-admin",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock, emailRepoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, emailRepoMock)

			err := uc.SendEmail(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.SendEmail() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_DeleteKolByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          uuid.UUID
	}{
		{
			name:          "success",
			wantErr:       false,
			expectedError: nil,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)
				repoMock.EXPECT().DeleteKolByID(gomock.Any(), gomock.Any()).Return(nil)

				return repoMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "delete_kol_by_id_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.DeleteKolByID error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)
				repoMock.EXPECT().DeleteKolByID(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

				return repoMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "delete_kol_by_id_not_found",
			wantErr:       true,
			expectedError: NotFoundError{resource: "kol", id: "0193487b-f1a2-7a72-8ae4-197b84dc52d6"},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)
				repoMock.EXPECT().DeleteKolByID(gomock.Any(), gomock.Any()).Return(domain.ErrDataNotFound)

				return repoMock
			},
			args: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock, nil)

			err := uc.DeleteKolByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.DeleteKolByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.DeleteKolByID() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
