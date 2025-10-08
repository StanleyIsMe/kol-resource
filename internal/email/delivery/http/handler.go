package http

import (
	"kolresource/internal/email/usecase"	

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailUsecase usecase.EmailUseCase
}

func NewEmailHandler(emailUsecase usecase.EmailUseCase) *EmailHandler {
	return &EmailHandler{emailUsecase: emailUsecase}
}

func (h *EmailHandler) CreateEmailSender(c *gin.Context) {

}

func (h *EmailHandler) ListEmailSenders(c *gin.Context) {

}

func (h *EmailHandler) SendEmail(c *gin.Context) {

}
