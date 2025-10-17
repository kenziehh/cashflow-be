package dto

import (
	"github.com/kenziehh/cashflow-be/internal/domain/transaction/entity"
)

type CreateTransactionRequest struct {
	TransactionType string  `json:"transaction_type" validate:"required,oneof=income expense"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	CategoryID      string  `json:"category_id" validate:"required,ulid" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Note            string  `json:"note,omitempty"`
	Period          string  `json:"period" validate:"required,oneof=daily weekly monthly yearly"`
	Date            string  `json:"date" validate:"required,datetime=2006-01-02"`
	ProofFile       string  `json:"proof_file,omitempty"`
}

type UpdateTransactionRequest struct {
	// TransactionId   uuid.UUID `json:"transaction_id" validate:"required,uuid4"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=income expense"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	CategoryID      string  `json:"category_id" validate:"required,ulid" swaggertype:"string" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Note            string  `json:"note,omitempty"`
	Period          string  `json:"period" validate:"required,oneof=daily weekly monthly yearly"`
	Date            string  `json:"date" validate:"required,datetime=2006-01-02"`
	ProofFile       string  `json:"proof_file,omitempty"`
}



type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	PageSize     int `json:"page_size"`
}

type TransactionListParams struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Type      string `query:"type"`
	Period    string `query:"period"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	SortBy    string `query:"sort_by"`
	OrderBy   string `query:"order_by"`
}

type PaginatedTransactionsResponse struct {
	Data        []*entity.Transaction `json:"data"`
	CurrentPage int                   `json:"current_page"`
	Limit       int                   `json:"limit"`
	TotalPage   int                   `json:"total_page"`
}

type SummaryTransactionResponse struct {
	TotalIncomeMonthly  float64 `json:"total_income_monthly"`
	TotalExpenseMonthly float64 `json:"total_expense_monthly"`
	TotalIncomeDaily    float64 `json:"total_income_daily"`
	TotalExpenseDaily   float64 `json:"total_expense_daily"`
}