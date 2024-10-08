package entities

import (
	"time"

	"github.com/google/uuid"
)

type SendEmailLog struct {
	ID          uuid.UUID `json:"id"`
	KolID       uuid.UUID `json:"kol_id"`
	KolName     string    `json:"kol_name"`
	Email       string    `json:"email"`
	AdminID     uuid.UUID `json:"admin_id"`
	AdminName   string    `json:"admin_name"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	CreatedAt   time.Time `json:"created_at"`
}
