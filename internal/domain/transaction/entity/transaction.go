package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	CategoryID      string    `json:"category_id" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	TransactionType string    `json:"transaction_type"` // e.g., "income" or "expense"
	Amount          float64   `json:"amount"`
	Period          string    `json:"period"`
	Note            string    `json:"note"`
	Date            string    `json:"date"`
	ProofFile       string   `json:"proof_file,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
