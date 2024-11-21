package sqlboiler

import (
	"context"
	"kolresource/internal/kol"
	"kolresource/internal/kol/domain"
	"kolresource/internal/kol/domain/entities"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreateKol(t *testing.T) {
	tests := []struct {
		name    string
		param   domain.CreateKolParams
		wantErr bool
	}{
		{
			name: "happy path",
			param: domain.CreateKolParams{
				Name:           "Test KOL",
				Email:          "test@example.com",
				SocialMedia:    "instagram",
				Description:    "Test description",
				Sex:            "m",
				Enable:         true,
				UpdatedAdminID: uuid.Must(uuid.NewV7()),
			},
			wantErr: false,
		},
		{
			name: "dumplicate email error",
			param: domain.CreateKolParams{
				Name:           "Test KOL",
				Email:          "test@example.com",
				Description:    "Test description",
				Sex:            "m",
				Enable:         true,
				UpdatedAdminID: uuid.Must(uuid.NewV7()),
			},
			wantErr: true,
		},
		{
			name: "very long name",
			param: domain.CreateKolParams{
				Name:           strings.Repeat("a", 256),
				Email:          "test@example.com",
				Description:    "Test description",
				Sex:            "f",
				Enable:         true,
				UpdatedAdminID: uuid.Must(uuid.NewV7()),
			},
			wantErr: true,
		},
		{
			name: "invalid sex",
			param: domain.CreateKolParams{
				Name:           "Test KOL",
				Email:          "test@example.com",
				Description:    "Test description",
				Sex:            "invalid",
				Enable:         true,
				UpdatedAdminID: uuid.Must(uuid.NewV7()),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			repo := NewKolRepository(suitUpInstance.stdConn)
			ctx := context.Background()

			kol, err := repo.CreateKol(ctx, tt.param)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateKol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if kol == nil {
					t.Error("CreateKol() returned nil KOL when no error was expected")
					return
				}

				// Verify the created KOL
				if kol.Name != tt.param.Name {
					t.Errorf("Created KOL name = %v, want %v", kol.Name, tt.param.Name)
				}
				if kol.Email != tt.param.Email {
					t.Errorf("Created KOL email = %v, want %v", kol.Email, tt.param.Email)
				}
				if kol.Description != tt.param.Description {
					t.Errorf("Created KOL description = %v, want %v", kol.Description, tt.param.Description)
				}
				if kol.SocialMedia != tt.param.SocialMedia {
					t.Errorf("Created KOL socialMedia = %v, want %v", kol.SocialMedia, tt.param.SocialMedia)
				}
				if kol.Sex != tt.param.Sex {
					t.Errorf("Created KOL sex = %v, want %v", kol.Sex, tt.param.Sex)
				}
				if kol.Enable != tt.param.Enable {
					t.Errorf("Created KOL enable = %v, want %v", kol.Enable, tt.param.Enable)
				}
				if kol.UpdatedAdminID != tt.param.UpdatedAdminID {
					t.Errorf("Created KOL updatedAdminID = %v, want %v", kol.UpdatedAdminID, tt.param.UpdatedAdminID)
				}
			}
		})
	}
}

func TestGetKolByID(t *testing.T) {
	repo := NewKolRepository(suitUpInstance.stdConn)
	ctx := context.Background()

	newKol, err := repo.CreateKol(ctx, domain.CreateKolParams{
		Name:           "stanley",
		Email:          "TestGetKolByID@gmail.com",
		SocialMedia:    "instagram",
		Description:    "description",
		Sex:            "m",
		Enable:         true,
		UpdatedAdminID: uuid.Must(uuid.NewV7()),
	})
	if err != nil {
		t.Fatalf("Failed to create test KOL: %v", err)
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		wantKol *entities.Kol
		wantErr bool
	}{
		{
			name:    "happy",
			id:      newKol.ID,
			wantKol: newKol,
			wantErr: false,
		},
		{
			name:    "data not found",
			id:      uuid.Must(uuid.NewV7()),
			wantKol: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			gotKol, err := repo.GetKolByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKolByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if gotKol == nil {
					t.Error("GetKolByID() returned nil KOL when no error was expected")
					return
				}

				if !reflect.DeepEqual(gotKol, tt.wantKol) {
					t.Errorf("got %v, but want %v", gotKol, tt.wantKol)
				}
			} else if gotKol != nil {
				t.Errorf("GetKolByID() = %v, want nil", gotKol)
			}
		})
	}
}

func TestGetKolByEmail(t *testing.T) {
	repo := NewKolRepository(suitUpInstance.stdConn)
	ctx := context.Background()

	newKol, err := repo.CreateKol(ctx, domain.CreateKolParams{
		Name:           "stanley",
		Email:          "TestGetKolByEmail@gmail.com",
		SocialMedia:    "instagram",
		Description:    "description",
		Sex:            "m",
		Enable:         true,
		UpdatedAdminID: uuid.Must(uuid.NewV7()),
	})
	if err != nil {
		t.Fatalf("Failed to create test KOL: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		wantKol *entities.Kol
		wantErr bool
	}{
		{
			name:    "happy",
			email:   newKol.Email,
			wantKol: newKol,
			wantErr: false,
		},
		{
			name:    "data not found",
			email:   "notfound@example.com",
			wantKol: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			gotKol, err := repo.GetKolByEmail(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKolByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if gotKol == nil {
					t.Error("GetKolByID() returned nil KOL when no error was expected")
					return
				}

				if !reflect.DeepEqual(gotKol, tt.wantKol) {
					t.Errorf("got %v, but want %v", gotKol, tt.wantKol)
				}
			} else if gotKol != nil {
				t.Errorf("GetKolByID() = %v, want nil", gotKol)
			}
		})
	}
}

func TestUpdateKol(t *testing.T) {
	ctx := context.Background()
	repo := NewKolRepository(suitUpInstance.stdConn)

	// Create a test KOL to update
	initialKol, err := repo.CreateKol(ctx, domain.CreateKolParams{
		Name:           "Initial Name",
		Email:          "initial@example.com",
		Description:    "Initial description",
		SocialMedia:    "instagram",
		Sex:            kol.SexMale,
		Enable:         true,
		UpdatedAdminID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("Failed to create initial test KOL: %v", err)
	}

	tests := []struct {
		name    string
		param   domain.UpdateKolParams
		wantKol *entities.Kol
		wantErr bool
	}{
		{
			name: "happy path",
			param: domain.UpdateKolParams{
				ID:             initialKol.ID,
				Name:           "Updated Name",
				Email:          "updated@example.com",
				SocialMedia:    "Updated instagram",
				Description:    "Updated description",
				Sex:            kol.SexFemale,
				Enable:         false,
				UpdatedAdminID: uuid.MustParse("6a9f3135-8415-4002-8011-f8c1d283964e"),
			},
			wantKol: &entities.Kol{
				ID:             initialKol.ID,
				Name:           "Updated Name",
				Email:          "updated@example.com",
				SocialMedia:    "Updated instagram",
				Description:    "Updated description",
				Sex:            kol.SexFemale,
				Enable:         false,
				UpdatedAdminID: uuid.MustParse("6a9f3135-8415-4002-8011-f8c1d283964e"),
				CreatedAt:      initialKol.CreatedAt,
			},
			wantErr: false,
		},
		{
			name: "non-existent KOL",
			param: domain.UpdateKolParams{
				ID:             uuid.New(),
				Name:           "Non-existent",
				Email:          "nonexistent@example.com",
				SocialMedia:    "instagram",
				Description:    "This KOL doesn't exist",
				Sex:            kol.SexMale,
				Enable:         true,
				UpdatedAdminID: uuid.New(),
			},
			wantKol: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKol, err := repo.UpdateKol(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateKol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(gotKol, tt.wantKol) {
					t.Errorf("UpdateKol() = %v, want %v", gotKol, tt.wantKol)
				}
			}
		})
	}
}
