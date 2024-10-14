package middleware

import (
	"kolresource/internal/admin"
	"kolresource/internal/admin/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthHeader  = "Authorization"
	BearerToken = "Bearer"
)

func JWT(usecase usecase.AdminUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(AuthHeader)
		if token == "" || !strings.HasPrefix(token, BearerToken) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

			return
		}

		token = strings.TrimSpace(strings.TrimPrefix(token, BearerToken))

		claims, err := usecase.LoginTokenParser(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

			return
		}

		c.Set(admin.AdminIDKey, claims.AdminID)
		c.Set(admin.AdminNameKey, claims.AdminName)

		c.Next()
	}
}
