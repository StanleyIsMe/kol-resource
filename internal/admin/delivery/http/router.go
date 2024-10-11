package http

import (
	"kolresource/internal/admin/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(router *gin.Engine, adminUsecase usecase.AdminUseCase) {
	adminHandler := NewAdminHandler(adminUsecase)

	v1 := router.Group("/api/v1")
	v1.POST("/register", adminHandler.Register)
	v1.POST("/login", adminHandler.Login)
}
