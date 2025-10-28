package entities

import (
	"time"

	"github.com/google/uuid"
)

const (
	DefaultDaRateLimit = 500
)

type EmailSender struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Key            string    `json:"key"`
	RateLimit      int       `json:"rate_limit"`
	Enabled        bool      `json:"enabled"`
	UpdatedAdminID uuid.UUID `json:"updated_admin_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastSendAt     time.Time `json:"last_send_at"`
}

func (e *EmailSender) SetDailyRateLimit(rateLimit int) {
	if rateLimit <= 0 {
		e.RateLimit = DefaultDaRateLimit

		return
	}

	e.RateLimit = rateLimit
}
