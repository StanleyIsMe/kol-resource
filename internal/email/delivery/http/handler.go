package http

import (
	"errors"
	"fmt"
	"kolresource/internal/email/schedule"
	"kolresource/internal/email/usecase"
	"kolresource/pkg/business"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type EmailHandler struct {
	emailUsecase  usecase.EmailUseCase
	emailSchedule *schedule.EmailSchedule
}

func NewEmailHandler(emailUsecase usecase.EmailUseCase, emailSchedule *schedule.EmailSchedule) *EmailHandler {
	return &EmailHandler{emailUsecase: emailUsecase, emailSchedule: emailSchedule}
}

// @Summary Create a new email sender
// @Description Create a new email sender for sending emails
// @Tags kol
// @Accept json
// @Produce json
// @Param request body CreateEmailSenderRequest true "Create email sender request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_senders [post]
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

// @Summary List email senders
// @Description List email senders
// @Tags email
// @Accept json
// @Produce json
// @Success 200 {object} ListEmailSendersResponse "Email senders list"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_senders [get]
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

// @Summary Update an email sender
// @Description Update an email sender
// @Tags email
// @Accept json
// @Produce json
// @Param id path string true "Email sender ID"
// @Param request body UpdateEmailSenderRequest true "Update email sender request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_senders/:id [put]
func (h *EmailHandler) UpdateEmailSender(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateEmailSenderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid uri param")})

		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	ucParam, err := req.ToUsecaseParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

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

// @Summary Get an email sender by id
// @Description Get an email sender by id
// @Tags email
// @Accept json
// @Produce json
// @Param id path string true "Email sender ID"
// @Success 200 {object} usecase.EmailSender "Email sender details"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_senders/:id [get]
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

// @Summary Send email
// @Description Send email
// @Tags email
// @Accept json
// @Produce json
// @Param request body SendEmailRequest true "Send email request"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_jobs [post]
func (h *EmailHandler) SendEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})

		return
	}

	if err := h.emailUsecase.SendEmail(ctx, req.ToUsecaseParam(c)); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("email send error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary List email jobs
// @Description List email jobs
// @Tags email
// @Accept json
// @Produce json
// @Success 200 {object} ListEmailJobsResponse "Email jobs list"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_jobs [get]
func (h *EmailHandler) ListEmailJobs(c *gin.Context) {
	ctx := c.Request.Context()

	var req ListEmailJobsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
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

// @Summary Get an email job by id
// @Description Get an email job details, including email logs
// @Tags email
// @Accept json
// @Produce json
// @Param id path string true "Email job ID"
// @Success 200 {object} GetEmailJobResponse "Email job details"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_jobs/:id [get]
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

// @Summary Cancel an email job
// @Description Cancel an email job, only pending or processing jobs can be canceled
// @Tags email
// @Accept json
// @Produce json
// @Param id path string true "Email job ID"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_jobs/:id/cancel [put]
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

// @Summary Start an email job
// @Description Start an email job, only canceled jobs can be started
// @Tags email
// @Accept json
// @Produce json
// @Param id path string true "Email job ID"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} nil "invalid request"
// @Failure 500 {object} business.ErrorResponse "internal error"
// @Router /api/v1/email_jobs/:id/start [put]
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

func (h *EmailHandler) EmailScheduleDebug(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.emailSchedule.SendEmailJob(ctx); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"error": err,
		}).Msg("email schedule send email error")

		c.JSON(business.UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
