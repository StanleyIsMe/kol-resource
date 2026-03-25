package http

import (
	"context"
	apiCfg "kolresource/internal/api/config"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/repository/email"
	"kolresource/internal/email/schedule"
	"kolresource/internal/email/usecase"
	"kolresource/pkg/config"

	"github.com/gin-gonic/gin"
)

type RegisterEmailRoutesParams struct {
	Router          *gin.RouterGroup
	Cfg             *config.Config[apiCfg.Config]
	EmailUsecase    usecase.EmailUseCase
	EmailRepository domain.Repository
}

func RegisterEmailRoutes(ctx context.Context, router *gin.RouterGroup, params RegisterEmailRoutesParams) {
	// sendEmailRepository := email.NewRepository(params.Cfg)
	sendEmailRepository := email.NewMockRepository(params.Cfg)
	emailSchedule := schedule.NewEmailSchedule(params.EmailRepository, sendEmailRepository, 0)
	emailHandler := NewEmailHandler(params.EmailUsecase, emailSchedule)
	emailSchedule.Start(ctx)
	v1 := router.Group("/api/v1")

	v1.POST("/email_senders", emailHandler.CreateEmailSender)
	v1.GET("/email_senders", emailHandler.ListEmailSenders)
	v1.PUT("/email_senders/:id", emailHandler.UpdateEmailSender)
	v1.GET("/email_senders/:id", emailHandler.GetEmailSender)

	v1.GET("/email_jobs", emailHandler.ListEmailJobs)
	v1.POST("/email_jobs", emailHandler.SendEmail)
	v1.GET("/email_jobs/:id", emailHandler.GetEmailJob)
	v1.PUT("/email_jobs/:id/cancel", emailHandler.CancelEmailJob)
	v1.PUT("/email_jobs/:id/start", emailHandler.StartEmailJob)

	v1.GET("/email_schedule", emailHandler.EmailScheduleDebug)
}
