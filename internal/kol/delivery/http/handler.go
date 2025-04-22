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

// @Summary Create a new kol
// @Description Create a new kol
// @Tags kol
// @Accept json
// @Produce json
// @Param request body CreateKolRequest true "Create kol request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/kols [post]
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

func (h *KolHandler) BatchCreateKolsByXlsx(c *gin.Context) {
	ctx := c.Request.Context()

	var req BatchCreateKolsByXlsxRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := h.uc.BatchCreateKolsByXlsx(ctx, req.ToUsecaseParam(c)); err != nil {
		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Update a kol
// @Description Update a kol
// @Tags kol
// @Accept json
// @Produce json
// @Param id path string true "Kol ID"
// @Param request body UpdateKolRequest true "Update kol request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/kols [put]
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

// @Summary Get a kol by id
// @Description Get a kol by id
// @Tags kol
// @Accept json
// @Produce json
// @Param id path string true "Kol ID"
// @Success 200 {object} usecase.Kol "Kol details"
// @Failure 400 {object} nil "invalid kol id"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/kols/{id} [get]
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

// @Summary List kols
// @Description List kols
// @Tags kol
// @Accept json
// @Produce json
// @Param request query ListKolsRequest true "List kols request"
// @Success 200 {object} ListKolsResponse "Kol list"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/kols [get]
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

// @Summary Create a new tag
// @Description Create a new tag
// @Tags tag
// @Accept json
// @Produce json
// @Param request body CreateTagRequest true "Create tag request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/tags [post]
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

// @Summary List tags
// @Description List tags
// @Tags tag
// @Accept json
// @Produce json
// @Param name query string false "Tag name"
// @Success 200 {object} []usecase.Tag "Tag list"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/tags [get]
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

// @Summary Create a new product
// @Description Create a new product
// @Tags product
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Create product request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/products [post]
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

// @Summary List products
// @Description List products
// @Tags product
// @Accept json
// @Produce json
// @Param name query string false "Product name"
// @Success 200 {object} []usecase.Product "Product list"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/products [get]
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

// @Summary Send email
// @Description Send email
// @Tags email
// @Accept json
// @Produce json
// @Param request body SendEmailRequest true "Send email request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email [post]
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
