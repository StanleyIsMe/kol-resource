package handler

import (
	"kolresource/internal/admin"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAdminIDFromContext(c *gin.Context) uuid.UUID {
	adminID, ok := c.Get(admin.AdminIDKey)
	if !ok {
		return uuid.Nil
	}

	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return adminUUID
}
