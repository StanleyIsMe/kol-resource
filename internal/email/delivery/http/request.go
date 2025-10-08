package http

import (
	"kolresource/internal/admin"
	"kolresource/internal/email/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendEmailRequest struct {
	Subject      string           `json:"subject" binding:"required,gte=1,lte=100"`
	EmailContent string           `json:"email_content" binding:"required,gte=1"`
	KolIDs       []uuid.UUID      `json:"kol_ids" binding:"required"`
	ProductID    uuid.UUID        `json:"product_id" binding:"required"`
	Images       []SendEmailImage `json:"images" binding:"dive"`
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
		UpdatedAdminID:   GetAdminIDFromContext(c),
		UpdatedAdminName: c.GetString(admin.AdminNameKey),
		ProductID:        r.ProductID,
		Images:           images,
	}
}

func GetAdminIDFromContext(c *gin.Context) uuid.UUID {
	adminID, ok := c.Get(admin.AdminIDKey)
	if !ok {
		return uuid.Nil
	}

	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return adminUUID
}
