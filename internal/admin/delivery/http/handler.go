package http

import (
	"errors"
	"fmt"
	"kolresource/internal/admin/usecase"
	"kolresource/pkg/business"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AdminHandler struct {
	adminUsecase usecase.AdminUseCase
}

func NewAdminHandler(adminUsecase usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUsecase: adminUsecase}
}

// @Summary Register admin
// @Description Register a new admin user
// @Tags admin
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/register [post]
func (h *AdminHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.adminUsecase.Register(ctx, usecase.RegisterParams{
		Name:     req.Name,
		UserName: req.UserName,
		Password: req.Password,
	}); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("admin register error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Login admin
// @Description Login with admin credentials
// @Tags admin
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	resp, err := h.adminUsecase.Login(ctx, req.UserName, req.Password)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("admin login error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:     resp.Token,
		AdminName: resp.AdminName,
	})
}
