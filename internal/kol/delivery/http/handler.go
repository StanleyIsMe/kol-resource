package http

import (
	"errors"
	"fmt"
	"kolresource/internal/kol/usecase"
	"kolresource/pkg/business"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type KolHandler struct {
	uc usecase.KolUseCase
}

func NewKolHandler(uc usecase.KolUseCase) *KolHandler {
	return &KolHandler{uc: uc}
}

func (h *KolHandler) CreateKol(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateKolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.uc.CreateKol(ctx, req.ToUsecaseParam(c)); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("kol create error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *KolHandler) UpdateKol(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateKolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	ucParam, err := req.ToUsecaseParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := h.uc.UpdateKol(ctx, ucParam); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("kol update error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *KolHandler) GetKolByID(c *gin.Context) {
	ctx := c.Request.Context()

	kolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid kol id")})

		return
	}

	kol, err := h.uc.GetKolByID(ctx, kolID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": kolID,
			"error":   err,
		}).Msg("kol get by id error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, kol)
}

func (h *KolHandler) ListKols(c *gin.Context) {
	ctx := c.Request.Context()

	var req ListKolsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	ucParam, err := req.ToUsecaseParam()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	kols, total, err := h.uc.ListKols(ctx, ucParam)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("kol search error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, ListKolsResponse{
		Kols:  kols,
		Total: total,
	})
}

func (h *KolHandler) CreateTag(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.uc.CreateTag(ctx, usecase.CreateTagParam{
		Name:           req.Name,
		UpdatedAdminID: uuid.Must(uuid.NewV7()),
	}); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("tag create error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *KolHandler) ListTags(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.Query("name")

	tags, err := h.uc.ListTagsByName(ctx, name)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": name,
			"error":   err,
		}).Msg("tag search error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *KolHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.uc.CreateProduct(ctx, req.ToUsecaseParam(c)); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("product create error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *KolHandler) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.Query("name")

	products, err := h.uc.ListProductsByName(ctx, name)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": name,
			"error":   err,
		}).Msg("product search error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *KolHandler) SendEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.uc.SendEmail(ctx, req.ToUsecaseParam(c)); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("email send error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
