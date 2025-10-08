package entities

import (
	"time"

	"kolresource/internal/email"

	"github.com/google/uuid"
)

type EmailLog struct {
	ID        int64                `json:"id"`
	JobID     int64                `json:"job_id"`
	KolID     uuid.UUID            `json:"kol_id"`
	KolName   string               `json:"kol_name"`
	Email     string               `json:"email"`
	Status    email.EmailLogStatus `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
