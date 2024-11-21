package email

import (
	"context"
	"fmt"
	apiCfg "kolresource/internal/api/config"
	"kolresource/internal/kol/domain"
	"kolresource/pkg/config"

	"gopkg.in/gomail.v2"
)

type Repository struct {
	cfg *config.Config[apiCfg.Config]
}

func NewRepository(cfg *config.Config[apiCfg.Config]) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (repo *Repository) SendEmail(_ context.Context, param domain.SendEmailParams) error {
	dialer := gomail.NewDialer(repo.cfg.CustomConfig.Email.ServerHost, repo.cfg.CustomConfig.Email.ServerPort, param.AdminEmail, param.AdminPass)
	sendCloser, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial mail server: %w", err)
	}

	mailMsg := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	for _, toEmail := range param.ToEmails {
		mailMsg.SetHeader("From", mailMsg.FormatAddress(param.AdminEmail, repo.cfg.CustomConfig.Email.AdminName))
		mailMsg.SetAddressHeader("To", toEmail.Email, toEmail.Name)
		mailMsg.SetHeader("Subject", param.Subject)
		mailMsg.SetHeader("To", toEmail.Email)
		mailContent := fmt.Sprintf("您好 %s, \n\n %s", toEmail.Name, param.Body)
		mailMsg.SetBody("text/html", mailContent)

		if err := gomail.Send(sendCloser, mailMsg); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		mailMsg.Reset()
	}

	return nil
}
