package dto

type MaximumSpendRequest struct {
	ID           string  `json:"id,omitempty"`
	DailyLimit   float64 `json:"daily_limit" validate:"gte=0"`
	MonthlyLimit float64 `json:"monthly_limit" validate:"gte=0"`
	YearlyLimit  float64 `json:"yearly_limit" validate:"gte=0"`
}

type MaximumSpendResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	DailyLimit   float64 `json:"daily_limit"`
	MonthlyLimit float64 `json:"monthly_limit"`
	YearlyLimit  float64 `json:"yearly_limit"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
