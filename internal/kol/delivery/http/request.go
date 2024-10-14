package http

import (
	"fmt"
	"kolresource/internal/admin"
	"kolresource/internal/kol"
	"kolresource/internal/kol/usecase"
	"kolresource/pkg/transport/pager"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateKolRequest struct {
	Name        string      `json:"name" binding:"required"`
	Email       string      `json:"email" binding:"required,email"`
	Description string      `json:"description" binding:"lte=500"`
	Sex         kol.Sex     `json:"sex" binding:"required,oneof=m f"`
	Tags        []uuid.UUID `json:"tags"`
}

func (r *CreateKolRequest) ToUsecaseParam(c *gin.Context) usecase.CreateKolParam {
	return usecase.CreateKolParam{
		Name:           r.Name,
		Email:          r.Email,
		Description:    r.Description,
		Sex:            r.Sex,
		Tags:           r.Tags,
		UpdatedAdminID: GetAdminIDFromContext(c),
	}
}

type UpdateKolRequest struct {
	Name        string      `json:"name" binding:"required"`
	Email       string      `json:"email" binding:"required,email"`
	Description string      `json:"description" binding:"lte=500"`
	Sex         kol.Sex     `json:"sex" binding:"required,oneof=m f"`
	Tags        []uuid.UUID `json:"tags"`
}

func (r *UpdateKolRequest) ToUsecaseParam(c *gin.Context) (usecase.UpdateKolParam, error) {
	kolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return usecase.UpdateKolParam{}, fmt.Errorf("invalid kol id: %w", err)
	}

	return usecase.UpdateKolParam{
		KolID:          kolID,
		Name:           r.Name,
		Email:          r.Email,
		Description:    r.Description,
		Sex:            r.Sex,
		Tags:           r.Tags,
		UpdatedAdminID: GetAdminIDFromContext(c),
	}, nil
}

type ListKolsRequest struct {
	Email *string  `form:"email,omitempty"`
	Name  *string  `form:"name,omitempty"`
	Tag   *string  `form:"tag,omitempty"`
	Sex   *kol.Sex `form:"sex,omitempty" binding:"omitempty,oneof=m f"`

	pager.Page
}

func (r *ListKolsRequest) ToUsecaseParam() usecase.ListKolsParam {
	return usecase.ListKolsParam{
		Email:    r.Email,
		Name:     r.Name,
		Tag:      r.Tag,
		Sex:      r.Sex,
		Page:     r.PageIndex,
		PageSize: r.PageSize,
	}
}

type ListKolsResponse struct {
	Kols  []*usecase.Kol `json:"kols"`
	Total int            `json:"total"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,lte=50"`
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required,lte=50"`
	Description string `json:"description" binding:"lte=500"`
}

func (r *CreateProductRequest) ToUsecaseParam(c *gin.Context) usecase.CreateProductParam {
	return usecase.CreateProductParam{
		Name:           r.Name,
		Description:    r.Description,
		UpdatedAdminID: GetAdminIDFromContext(c),
	}
}

type SendEmailRequest struct {
	Subject      string      `json:"subject" binding:"required,gte=1,lte=100"`
	EmailContent string      `json:"email_content" binding:"required,gte=1"`
	KolIDs       []uuid.UUID `json:"kol_ids" binding:"required"`
	ProductID    uuid.UUID   `json:"product_id" binding:"required"`
}

func (r *SendEmailRequest) ToUsecaseParam(c *gin.Context) usecase.SendEmailParam {
	return usecase.SendEmailParam{
		Subject:          r.Subject,
		EmailContent:     r.EmailContent,
		KolIDs:           r.KolIDs,
		UpdatedAdminID:   GetAdminIDFromContext(c),
		UpdatedAdminName: c.GetString(admin.AdminNameKey),
		ProductID:        r.ProductID,
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