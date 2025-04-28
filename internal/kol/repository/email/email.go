package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	apiCfg "kolresource/internal/api/config"
	"kolresource/internal/kol/domain"
	"kolresource/pkg/config"

	"github.com/rs/zerolog"
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

type MailImage struct {
	Filename string
	Data     []byte
	Header   map[string][]string
}

func (repo *Repository) SendEmail(ctx context.Context, param domain.SendEmailParams) error {
	dialer := gomail.NewDialer(
		repo.cfg.CustomConfig.Email.ServerHost,
		repo.cfg.CustomConfig.Email.ServerPort,
		repo.cfg.CustomConfig.Email.AdminEmail,
		repo.cfg.CustomConfig.Email.AdminPass,
	)
	sendCloser, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial mail server: %w", err)
	}

	emailImages := make([]MailImage, 0, len(param.Images))
	for _, img := range param.Images {
		// Decode the base64 image data
		imgData, err := base64.StdEncoding.DecodeString(img.Data)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}

		// Determine content type based on image type
		contentType := fmt.Sprintf("image/%s", img.ImageType)

		// Create an inline attachment with the correct Content-ID
		// Note: The Content-ID in the email must match the 'cid:...' reference in the HTML
		contentID := fmt.Sprintf("<%s>", img.ContentID) // Format as <image1>

		emailImages = append(emailImages, MailImage{
			Filename: fmt.Sprintf("%s.%s", img.ContentID, img.ImageType),
			Data:     imgData,
			Header: map[string][]string{
				"Content-Type":        {contentType},
				"Content-ID":          {contentID},
				"Content-Disposition": {"inline"},
			},
		})
	}

	mailMsg := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	for _, toEmail := range param.ToEmails {
		mailMsg.SetHeader("From", mailMsg.FormatAddress(repo.cfg.CustomConfig.Email.AdminEmail, repo.cfg.CustomConfig.Email.AdminName))
		mailMsg.SetAddressHeader("To", toEmail.Email, toEmail.Name)
		mailMsg.SetHeader("Subject", param.Subject)
		mailMsg.SetHeader("To", toEmail.Email)
		mailContent := fmt.Sprintf("您好 %s, \n\n %s", toEmail.Name, param.Body)
		mailMsg.SetBody("text/html", mailContent)

		// Attach each image with its Content-ID
		for _, img := range emailImages {
			// Add the inline attachment
			mailMsg.Attach(
				img.Filename,
				// Set the content with the image data
				gomail.SetCopyFunc(func(w io.Writer) error {
					_, err := w.Write(img.Data)

					return err
				}),
				gomail.SetHeader(img.Header),
			)
		}

		if err := gomail.Send(sendCloser, mailMsg); err != nil {
			zerolog.Ctx(ctx).Error().Fields(map[string]any{
				"error":   err,
				"toEmail": toEmail,
			}).Msg("failed to send email")

			return fmt.Errorf("failed to send email: %w", err)
		}

		mailMsg.Reset()
	}

	return nil
}
