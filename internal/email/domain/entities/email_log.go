package entities

import (
	"time"

	"kolresource/internal/email"

	"github.com/google/uuid"
)

type EmailLog struct {
	ID        int64                `json:"id"`
	JobID     int64                `json:"job_id"`
	ProductID uuid.UUID            `json:"product_id"`
	SenderID  uuid.UUID            `json:"sender_id"`
	MessageID string               `json:"message_id"`
	KolID     uuid.UUID            `json:"kol_id"`
	KolName   string               `json:"kol_name"`
	Email     string               `json:"email"`
	Status    email.EmailLogStatus `json:"status"`
	Memo      string               `json:"memo"`
	Reply     bool                 `json:"reply"`
	CreatedAt time.Time            `json:"created_at"`
	SendedAt  time.Time            `json:"sended_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
