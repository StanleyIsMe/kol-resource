package http

import (
	"kolresource/internal/kol/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterKolRoutes(router *gin.RouterGroup, kolUsecase usecase.KolUseCase) {
	kolHandler := NewKolHandler(kolUsecase)

	v1 := router.Group("/api/v1")
	v1.POST("/kols", kolHandler.CreateKol)
	v1.PUT("/kols/:id", kolHandler.UpdateKol)
	v1.GET("/kols/:id", kolHandler.GetKolByID)
	v1.GET("/kols", kolHandler.ListKols)
	v1.POST("/kols/upload", kolHandler.BatchCreateKolsByXlsx)

	v1.POST("/tags", kolHandler.CreateTag)
	v1.GET("/tags", kolHandler.ListTags)

	v1.POST("/products", kolHandler.CreateProduct)
	v1.GET("/products", kolHandler.ListProducts)

	v1.POST("/send_emails", kolHandler.SendEmail)
}
