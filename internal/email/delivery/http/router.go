package http

import (
	"kolresource/internal/email/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterEmailRoutes(router *gin.RouterGroup, emailUsecase usecase.EmailUseCase) {
	emailHandler := NewEmailHandler(emailUsecase)

	v1 := router.Group("/api/v1")
	v1.POST("/send_emails", emailHandler.SendEmail)
}