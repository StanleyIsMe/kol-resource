package entities

import (
	"kolresource/internal/kol/domain"
	"time"

	"github.com/google/uuid"
)

type Kol struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Description    string     `json:"description"`
	Sex            domain.Sex `json:"sex"`
	Enable         bool       `json:"enable"`
	UpdatedAdminID uuid.UUID  `json:"updated_admin_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
