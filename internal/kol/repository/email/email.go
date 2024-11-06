package email

import (
	"context"
	"fmt"
	apiCfg "kolresource/internal/api/config"
	"kolresource/internal/kol/domain"
	"kolresource/pkg/config"

	"gopkg.in/gomail.v2"
)

type EmailRepository struct {
	cfg *config.Config[apiCfg.Config]
}

func NewEmailRepository(cfg *config.Config[apiCfg.Config]) *EmailRepository {
	return &EmailRepository{
		cfg: cfg,
	}
}

func (repo *EmailRepository) SendEmail(ctx context.Context, param domain.SendEmailParams) error {
	dialer := gomail.NewDialer(repo.cfg.CustomConfig.Email.ServerHost, repo.cfg.CustomConfig.Email.ServerPort, param.AdminEmail, param.AdminPass)
	sendCloser, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial mail server: %w", err)
	}

	mailMsg := gomail.NewMessage()
	for _, toEmail := range param.ToEmails {
		mailMsg.SetHeader("From", mailMsg.FormatAddress(param.AdminEmail, repo.cfg.CustomConfig.Email.AdminName))
		mailMsg.SetAddressHeader("To", toEmail.Email, toEmail.Name)
		mailMsg.SetHeader("Subject", param.Subject)
		mailMsg.SetHeader("To", toEmail.Email)
		mailContent := fmt.Sprintf("您好 %s, \n\n %s", toEmail.Name, param.Body)
		mailMsg.SetBody("text/html", mailContent)

		if err := gomail.Send(sendCloser, mailMsg); err != nil {
			// TODO: handle error
			fmt.Printf("failed to send email: %v", err)
		}

		mailMsg.Reset()
	}

	return nil
}
