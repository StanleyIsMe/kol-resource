package http

import (
	"kolresource/internal/email/usecase"

	"github.com/gin-gonic/gin"
)

type RegisterEmailRoutesParams struct {
	Router *gin.RouterGroup
	EmailUsecase usecase.EmailUseCase
}

func RegisterEmailRoutes(router *gin.RouterGroup, emailUsecase usecase.EmailUseCase) {
	emailHandler := NewEmailHandler(emailUsecase)

	v1 := router.Group("/api/v1")
	v1.POST("/send_emails", emailHandler.SendEmail)

	v1.POST("/email_senders", emailHandler.CreateEmailSender)
	v1.GET("/email_senders", emailHandler.ListEmailSenders)
	v1.PUT("/email_senders/:id", emailHandler.UpdateEmailSender)
	v1.GET("/email_senders/:id", emailHandler.GetEmailSender)

	v1.GET("/email_jobs", emailHandler.ListEmailJobs)
	v1.GET("/email_jobs/:id", emailHandler.GetEmailJob)
	v1.PUT("/email_jobs/:id/cancel", emailHandler.CancelEmailJob)
	v1.PUT("/email_jobs/:id/start", emailHandler.StartEmailJob)
}
