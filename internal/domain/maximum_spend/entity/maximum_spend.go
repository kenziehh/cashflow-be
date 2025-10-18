package entity

import (
	"time"

	"github.com/google/uuid"
)

type MaximumSpend struct {
	ID           string    `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	DailyLimit   float64   `json:"daily_limit" db:"daily_limit"`
	MonthlyLimit float64   `json:"monthly_limit" db:"monthly_limit"`
	YearlyLimit  float64   `json:"yearly_limit" db:"yearly_limit"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
