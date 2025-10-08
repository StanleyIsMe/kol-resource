package entities

import (
	"time"

	"kolresource/internal/email"

	"github.com/google/uuid"
)

type EmailJob struct {
	ID                   int64                    `json:"id"`
	ExpectedReciverCount int                      `json:"expected_reciver_count"`
	SuccessCount         int                      `json:"success_count"`
	AdminID              uuid.UUID                `json:"admin_id"`
	AdminName            string                   `json:"admin_name"`
	ProductID            uuid.UUID                `json:"product_id"`
	ProductName          string                   `json:"product_name"`
	SenderID             uuid.UUID                `json:"sender_id"`
	SenderName           string                   `json:"sender_name"`
	SenderEmail          string                   `json:"sender_email"`
	Memo                 string                   `json:"memo"`
	Payload              string                   `json:"payload"`
	Status               email.EmailJobStatus     `json:"status"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
	LastExecuteAt        time.Time                `json:"last_execute_at"`
}
