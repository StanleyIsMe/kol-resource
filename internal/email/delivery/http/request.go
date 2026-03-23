package http

import (
	"fmt"
	"kolresource/internal/admin"
	"kolresource/internal/common/handler"
	"kolresource/internal/email"
	"kolresource/internal/email/usecase"
	"kolresource/pkg/transport/pager"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendEmailRequest struct {
	Subject      string           `json:"subject" binding:"required,gte=1,lte=100"`
	EmailContent string           `json:"email_content" binding:"required,gte=1"`
	KolIDs       []uuid.UUID      `json:"kol_ids" binding:"required"`
	ProductID    uuid.UUID        `json:"product_id" binding:"required"`
	Images       []SendEmailImage `json:"images" binding:"dive"`
	SenderID     uuid.UUID        `json:"sender_id" binding:"required"`
}

type SendEmailImage struct {
	ContentID string `json:"content_id" binding:"required"`
	Data      string `json:"data" binding:"required"`
	ImageType string `json:"type" binding:"required"`
}

func (r *SendEmailRequest) ToUsecaseParam(c *gin.Context) usecase.SendEmailParam {
	images := make([]usecase.SendEmailImage, len(r.Images))
	for index, image := range r.Images {
		images[index] = usecase.SendEmailImage{
			ContentID: image.ContentID,
			Data:      image.Data,
			ImageType: image.ImageType,
		}
	}

	return usecase.SendEmailParam{
		Subject:          r.Subject,
		EmailContent:     r.EmailContent,
		KolIDs:           r.KolIDs,
		UpdatedAdminID:   handler.GetAdminIDFromContext(c),
		UpdatedAdminName: c.GetString(admin.AdminNameKey),
		ProductID:        r.ProductID,
		Images:           images,
		SenderID:         r.SenderID,
	}
}

type CreateEmailSenderRequest struct {
	Name      string `json:"name" binding:"required,lte=50"`
	Email     string `json:"email" binding:"required,email"`
	Key       string `json:"key" binding:"required"`
	RateLimit int    `json:"rate_limit" binding:"required,min=1"`
}

func (r *CreateEmailSenderRequest) ToUsecaseParam(c *gin.Context) usecase.CreateEmailSenderParam {
	return usecase.CreateEmailSenderParam{
		UpdatedAdminID:   handler.GetAdminIDFromContext(c),
		UpdatedAdminName: c.GetString(admin.AdminNameKey),
		Name:             r.Name,
		Email:            r.Email,
		Key:              r.Key,
		RateLimit:        r.RateLimit,
	}
}

type UpdateEmailSenderRequest struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	Name      *string `json:"name,omitempty" binding:"omitempty,lte=50"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Key       *string `json:"key,omitempty" binding:"omitempty,lte=255"`
	RateLimit *int    `json:"rate_limit,omitempty" binding:"omitempty,min=1"`
}

func (r *UpdateEmailSenderRequest) ToUsecaseParam(c *gin.Context) (usecase.UpdateEmailSenderParam, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return usecase.UpdateEmailSenderParam{}, fmt.Errorf("invalid uuid: %w", err)
	}

	return usecase.UpdateEmailSenderParam{
		ID:             id,
		Name:           r.Name,
		Email:          r.Email,
		Key:            r.Key,
		RateLimit:      r.RateLimit,
		UpdatedAdminID: handler.GetAdminIDFromContext(c),
	}, nil
}

type ListEmailSendersResponse struct {
	EmailSenders []usecase.EmailSender `json:"email_senders"`
	Total        int                   `json:"total"`
}

type ListEmailJobsRequest struct {
	SenderEmail *string          `form:"sender_email,omitempty"`
	SenderName  *string          `form:"sender_name,omitempty"`
	ProductName *string          `form:"product_name,omitempty"`
	Status      *email.JobStatus `form:"status,omitempty"`
	pager.Page
}

func (r *ListEmailJobsRequest) ToUsecaseParam() usecase.ListEmailJobsParam {
	return usecase.ListEmailJobsParam{
		SenderEmail: r.SenderEmail,
		SenderName:  r.SenderName,
		ProductName: r.ProductName,
		Status:      r.Status,
		Page:        r.PageIndex,
		PageSize:    r.PageSize,
	}
}

type GetEmailJobResponse struct {
	EmailJob  usecase.EmailJob   `json:"email_job"`
	EmailLogs []usecase.EmailLog `json:"email_logs"`
}
