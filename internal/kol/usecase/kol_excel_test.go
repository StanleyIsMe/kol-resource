package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain"
	"kolresource/internal/kol/domain/entities"
	"kolresource/internal/kol/mock/repositorymock"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func createTestXlsxFileHeader(t *testing.T, rows [][]string) *multipart.FileHeader {
	t.Helper()

	f := excelize.NewFile()
	sheetName := f.GetSheetList()[0]

	for i, row := range rows {
		for j, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				t.Fatalf("CoordinatesToCellName error: %v", err)
			}

			if err := f.SetCellValue(sheetName, cellName, cell); err != nil {
				t.Fatalf("SetCellValue error: %v", err)
			}
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("WriteToBuffer error: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.xlsx"`)
	h.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("CreatePart error: %v", err)
	}

	if _, err := part.Write(buf.Bytes()); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}

	reader := multipart.NewReader(body, writer.Boundary())

	form, err := reader.ReadForm(int64(body.Len()))
	if err != nil {
		t.Fatalf("ReadForm error: %v", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("no file found in form")
	}

	return files[0]
}

func TestKolUseCaseImpl_BatchCreateKolsByXlsx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		getFileHeader func(t *testing.T) *multipart.FileHeader
		adminID       uuid.UUID
	}{
		{
			name:    "success_create_new_kol",
			wantErr: false,
			getFileHeader: func(t *testing.T) *multipart.FileHeader {
				t.Helper()

				return createTestXlsxFileHeader(t, [][]string{
					{"tag1", "test-kol", "m", "https://social.com", "test@example.com"},
				})
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				tagID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

				repoMock.EXPECT().GetTagByName(gomock.Any(), "tag1").Return(&entities.Tag{
					ID:   tagID,
					Name: "tag1",
				}, nil)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "test@example.com").Return(nil, domain.ErrDataNotFound)

				repoMock.EXPECT().CreateKol(gomock.Any(), domain.CreateKolParams{
					Email:          "test@example.com",
					Name:           "test-kol",
					Sex:            kol.SexMale,
					SocialMedia:    "https://social.com",
					Enable:         true,
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Tags:           []uuid.UUID{tagID},
				}).Return(&entities.Kol{}, nil)

				return repoMock
			},
			adminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:    "success_update_existing_kol",
			wantErr: false,
			getFileHeader: func(t *testing.T) *multipart.FileHeader {
				t.Helper()

				return createTestXlsxFileHeader(t, [][]string{
					{"tag1", "test-kol", "f", "https://social.com", "existing@example.com"},
				})
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				tagID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
				kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

				repoMock.EXPECT().GetTagByName(gomock.Any(), "tag1").Return(&entities.Tag{
					ID:   tagID,
					Name: "tag1",
				}, nil)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "existing@example.com").Return(&entities.Kol{
					ID:    kolID,
					Email: "existing@example.com",
				}, nil)

				repoMock.EXPECT().UpdateKol(gomock.Any(), domain.UpdateKolParams{
					ID:             kolID,
					Email:          "existing@example.com",
					Name:           "test-kol",
					Sex:            kol.SexFemale,
					SocialMedia:    "https://social.com",
					Enable:         true,
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Tags:           []uuid.UUID{tagID},
				}).Return(&entities.Kol{}, nil)

				return repoMock
			},
			adminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:    "success_create_new_tag",
			wantErr: false,
			getFileHeader: func(t *testing.T) *multipart.FileHeader {
				t.Helper()

				return createTestXlsxFileHeader(t, [][]string{
					{"newtag", "test-kol", "m", "https://social.com", "test@example.com"},
				})
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				newTagID := uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6")

				repoMock.EXPECT().GetTagByName(gomock.Any(), "newtag").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateTag(gomock.Any(), domain.CreateTagParams{
					Name:           "newtag",
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				}).Return(&entities.Tag{ID: newTagID, Name: "newtag"}, nil)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "test@example.com").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateKol(gomock.Any(), domain.CreateKolParams{
					Email:          "test@example.com",
					Name:           "test-kol",
					Sex:            kol.SexMale,
					SocialMedia:    "https://social.com",
					Enable:         true,
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Tags:           []uuid.UUID{newTagID},
				}).Return(&entities.Kol{}, nil)

				return repoMock
			},
			adminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "GetTagByName_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetTagByName error: %w", errors.New("database error")),
			getFileHeader: func(t *testing.T) *multipart.FileHeader {
				t.Helper()

				return createTestXlsxFileHeader(t, [][]string{
					{"tag1", "test-kol", "m", "https://social.com", "test@example.com"},
				})
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetTagByName(gomock.Any(), "tag1").Return(nil, errors.New("database error"))

				return repoMock
			},
			adminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:          "CreateTag_error",
			wantErr:       true,
			expectedError: fmt.Errorf("uc.repo.CreateTag error: %w", errors.New("database error")),
			getFileHeader: func(t *testing.T) *multipart.FileHeader {
				t.Helper()

				return createTestXlsxFileHeader(t, [][]string{
					{"newtag", "test-kol", "m", "https://social.com", "test@example.com"},
				})
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetTagByName(gomock.Any(), "newtag").Return(nil, domain.ErrDataNotFound)
				repoMock.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			adminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
		{
			name:    "skip_invalid_row",
			wantErr: false,
			getFileHeader: func(t *testing.T) *multipart.FileHeader {
				t.Helper()

				return createTestXlsxFileHeader(t, [][]string{
					{"tag1", "test-kol", "m", "https://social.com", "invalid-email"},
				})
			},
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				return repoMock
			},
			adminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock)

			fileHeader := tt.getFileHeader(t)
			err := uc.BatchCreateKolsByXlsx(context.Background(), BatchCreateKolsByXlsxParam{
				File:           fileHeader,
				UpdatedAdminID: tt.adminID,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.BatchCreateKolsByXlsx() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.BatchCreateKolsByXlsx() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestKolUseCaseImpl_upsertKol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wantErr       bool
		expectedError error
		getRepoMock   func(ctrl *gomock.Controller) domain.Repository
		args          UpsertKolParam
	}{
		{
			name:    "create_new_kol",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "new@example.com").Return(nil, domain.ErrDataNotFound)

				repoMock.EXPECT().CreateKol(gomock.Any(), domain.CreateKolParams{
					Email:          "new@example.com",
					Name:           "new-kol",
					Sex:            kol.SexMale,
					SocialMedia:    "https://social.com",
					Enable:         true,
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Tags: []uuid.UUID{
						uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					},
				}).Return(&entities.Kol{}, nil)

				return repoMock
			},
			args: UpsertKolParam{
				Email:          "new@example.com",
				Name:           "new-kol",
				Sex:            kol.SexMale,
				SocialMedia:    "https://social.com",
				UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Tags: []uuid.UUID{
					uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				},
			},
		},
		{
			name:    "update_existing_kol",
			wantErr: false,
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "existing@example.com").Return(&entities.Kol{
					ID:    kolID,
					Email: "existing@example.com",
				}, nil)

				repoMock.EXPECT().UpdateKol(gomock.Any(), domain.UpdateKolParams{
					ID:             kolID,
					Email:          "existing@example.com",
					Name:           "updated-kol",
					Sex:            kol.SexFemale,
					SocialMedia:    "https://updated.com",
					Enable:         true,
					UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					Tags: []uuid.UUID{
						uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
					},
				}).Return(&entities.Kol{}, nil)

				return repoMock
			},
			args: UpsertKolParam{
				Email:          "existing@example.com",
				Name:           "updated-kol",
				Sex:            kol.SexFemale,
				SocialMedia:    "https://updated.com",
				UpdatedAdminID: uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				Tags: []uuid.UUID{
					uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				},
			},
		},
		{
			name:          "GetKolByEmail_error",
			wantErr:       true,
			expectedError: fmt.Errorf("repo.GetKolByEmail error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "test@example.com").Return(nil, errors.New("database error"))

				return repoMock
			},
			args: UpsertKolParam{
				Email: "test@example.com",
			},
		},
		{
			name:          "UpdateKol_error",
			wantErr:       true,
			expectedError: fmt.Errorf("uc.repo.UpdateKol error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "existing@example.com").Return(&entities.Kol{
					ID:    kolID,
					Email: "existing@example.com",
				}, nil)

				repoMock.EXPECT().UpdateKol(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: UpsertKolParam{
				Email: "existing@example.com",
			},
		},
		{
			name:          "CreateKol_error",
			wantErr:       true,
			expectedError: fmt.Errorf("uc.repo.CreateKol error: %w", errors.New("database error")),
			getRepoMock: func(ctrl *gomock.Controller) domain.Repository {
				repoMock := repositorymock.NewMockRepository(ctrl)

				repoMock.EXPECT().GetKolByEmail(gomock.Any(), "new@example.com").Return(nil, domain.ErrDataNotFound)

				repoMock.EXPECT().CreateKol(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repoMock
			},
			args: UpsertKolParam{
				Email: "new@example.com",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := tt.getRepoMock(ctrl)
			uc := NewKolUseCaseImpl(repoMock)

			err := uc.upsertKol(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("KolUseCaseImpl.upsertKol() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.expectedError.Error() {
				t.Errorf("KolUseCaseImpl.upsertKol() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
