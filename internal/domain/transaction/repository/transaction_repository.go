package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kenziehh/cashflow-be/internal/domain/transaction/dto"
	"github.com/kenziehh/cashflow-be/internal/domain/transaction/entity"
	"github.com/kenziehh/cashflow-be/pkg/errx"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *entity.Transaction) error
	GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, tx *entity.Transaction) error
	DeleteTransaction(ctx context.Context, id string) error
	GetTransactionsWithPagination(ctx context.Context, userID uuid.UUID, filter dto.TransactionListParams) (dto.PaginatedTransactionsResponse, error)
}

type transactionRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTransactionRepository(db *sql.DB, redis *redis.Client) TransactionRepository {
	return &transactionRepository{
		db:    db,
		redis: redis,
	}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, tx *entity.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, amount, type, category_id, note, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		uuid.New(),
		tx.UserID,
		tx.Amount,
		tx.TransactionType,
		tx.CategoryID,
		tx.Note,
		tx.Date,
		tx.CreatedAt,
		tx.UpdatedAt,
	)

	if err != nil {
		return errx.ErrDatabaseError
	}

	return nil
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, id string) (*entity.Transaction, error) {
	query := `
		SELECT id, user_id, amount, type, category_id, note, date, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	tx := &entity.Transaction{}
	err := row.Scan(
		&tx.ID,
		&tx.UserID,
		&tx.Amount,
		&tx.TransactionType,
		&tx.CategoryID,
		&tx.Note,
		&tx.Date,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errx.ErrTransactionNotFound
		}
		return nil, errx.ErrDatabaseError
	}

	return tx, nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, tx *entity.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, type = $2, category_id = $3, note = $4, date = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		tx.Amount,
		tx.TransactionType,
		tx.CategoryID,
		tx.Note,
		tx.Date,
		tx.UpdatedAt,
		tx.ID,
	)

	if err != nil {
		return errx.ErrDatabaseError
	}

	return nil
}

func (r *transactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	query := `
		DELETE FROM transactions
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errx.ErrDatabaseError
	}

	return nil
}

func (r *transactionRepository) GetTransactionsWithPagination(
	ctx context.Context,
	userID uuid.UUID,
	filter dto.TransactionListParams,
) (dto.PaginatedTransactionsResponse, error) {

	// fmt.Println("Filter received in repository:", filter)

	query := `
		SELECT id, user_id, amount, type, category_id, note, date, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	paramIndex := 2

	// Filter tanggal
	if filter.StartDate != "" && filter.EndDate != "" {
		query += fmt.Sprintf(" AND date >= $%d AND date <= $%d", paramIndex, paramIndex+1)
		args = append(args, filter.StartDate, filter.EndDate)
		paramIndex += 2
	}

	// Filter type
	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", paramIndex)
		args = append(args, filter.Type)
		paramIndex++
	}


	// Sort
	validSortColumns := map[string]bool{
		"date":       true,
		"amount":     true,
		"created_at": true,
	}
	if !validSortColumns[filter.SortBy] {
		filter.SortBy = "date"
	}

	order := strings.ToUpper(filter.OrderBy)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// Pagination
	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", filter.SortBy, order, paramIndex, paramIndex+1)
	args = append(args, filter.Limit, offset)

	// Eksekusi query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Println("Query error:", err)
		return dto.PaginatedTransactionsResponse{}, errx.ErrDatabaseError
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		tx := &entity.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Amount,
			&tx.TransactionType,
			&tx.CategoryID,
			&tx.Note,
			&tx.Date,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return dto.PaginatedTransactionsResponse{}, errx.ErrDatabaseError
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return dto.PaginatedTransactionsResponse{}, errx.ErrDatabaseError
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return dto.PaginatedTransactionsResponse{}, errx.ErrDatabaseError
	}

	response := dto.PaginatedTransactionsResponse{
		Data:        transactions,
		CurrentPage: filter.Page,
		Limit:       filter.Limit,
		TotalPage:   (total + filter.Limit - 1) / filter.Limit,
	}

	return response, nil
}
