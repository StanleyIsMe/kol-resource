package http

import (
	"errors"
	"fmt"
	"kolresource/internal/email/usecase"
	"kolresource/pkg/business"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type EmailHandler struct {
	emailUsecase usecase.EmailUseCase
}

func NewEmailHandler(emailUsecase usecase.EmailUseCase) *EmailHandler {
	return &EmailHandler{emailUsecase: emailUsecase}
}

func (h *EmailHandler) CreateEmailSender(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateEmailSenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	ucParam := req.ToUsecaseParam(c)
	if err := h.emailUsecase.CreateEmailSender(ctx, ucParam); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("email sender create error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *EmailHandler) ListEmailSenders(c *gin.Context) {
	ctx := c.Request.Context()

	emailSenders, err := h.emailUsecase.ListEmailSenders(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email sender list error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, ListEmailSendersResponse{
		EmailSenders: emailSenders,
		Total:        len(emailSenders),
	})
}

func (h *EmailHandler) UpdateEmailSender(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateEmailSenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	ucParam := req.ToUsecaseParam(c, id)
	if err := h.emailUsecase.UpdateEmailSender(ctx, ucParam); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("email sender update error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *EmailHandler) GetEmailSender(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	emailSender, err := h.emailUsecase.GetEmailSender(ctx, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email sender get error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, emailSender)
}

func (h *EmailHandler) SendEmail(_ *gin.Context) {

}

func (h *EmailHandler) ListEmailJobs(c *gin.Context) {
	ctx := c.Request.Context()

	var req ListEmailJobsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	emailJobs, err := h.emailUsecase.ListEmailJobs(ctx, req.ToUsecaseParam())
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email job list error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, emailJobs)
}

func (h *EmailHandler) GetEmailJob(c *gin.Context) {
	ctx := c.Request.Context()

	jobID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	emailJob, err := h.emailUsecase.GetEmailJob(ctx, jobID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email job get error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	emailLogs, err := h.emailUsecase.ListEmailLogs(ctx, usecase.ListEmailLogsParam{
		JobID: jobID,
	})
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email job get error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, GetEmailJobResponse{
		EmailJob:  *emailJob,
		EmailLogs: emailLogs,
	})
}

func (h *EmailHandler) CancelEmailJob(c *gin.Context) {
	ctx := c.Request.Context()

	jobID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.emailUsecase.CancelEmailJob(ctx, jobID); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email job cancel error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *EmailHandler) StartEmailJob(c *gin.Context) {
	ctx := c.Request.Context()

	jobID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.emailUsecase.StartEmailJob(ctx, jobID); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email job start error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
