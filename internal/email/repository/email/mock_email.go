package email

import (
	"context"
	apiCfg "kolresource/internal/api/config"
	"kolresource/internal/email/domain"
	"kolresource/pkg/config"
)

var _ domain.EmailRepository = (*MockRepository)(nil)

type MockRepository struct {
	cfg *config.Config[apiCfg.Config]
}

func NewMockRepository(cfg *config.Config[apiCfg.Config]) *MockRepository {
	return &MockRepository{
		cfg: cfg,
	}
}

func (repo *MockRepository) SendEmail(_ context.Context, _ domain.SendEmailParams) error {
	return nil
}
