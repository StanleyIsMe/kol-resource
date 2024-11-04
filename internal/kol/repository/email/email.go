package email

import (
	"context"
	"fmt"
	"kolresource/internal/kol/domain"

	"gopkg.in/gomail.v2"
)

type EmailRepository struct {
	mailServerHost string
	mailServerPort int
}

func NewEmailRepository(mailServerHost string, mailServerPort int) *EmailRepository {
	return &EmailRepository{
		mailServerHost: mailServerHost,
		mailServerPort: mailServerPort,
	}
}

func (repo *EmailRepository) SendEmail(ctx context.Context, param domain.SendEmailParams) error {
	dialer := gomail.NewDialer(repo.mailServerHost, repo.mailServerPort, param.AdminEmail, param.AdminPass)
	sendCloser, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial mail server: %w", err)
	}

	mailMsg := gomail.NewMessage()
	for _, toEmail := range param.ToEmails {
		mailMsg.SetHeader("From", param.AdminEmail)
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
