package repository

import (
	"context"
	"database/sql"
	"github.com/kenziehh/cashflow-be/pkg/errx"
	"github.com/kenziehh/cashflow-be/config/id"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kenziehh/cashflow-be/internal/domain/maximum_spend/entity"
)

type MaximumSpendRepository interface {
	CheckAlert(ctx context.Context, tx *entity.MaximumSpend, period string) error
	UpsertMaximumSpend(ctx context.Context, ms *entity.MaximumSpend) error
	GetMaximumSpendByUserID(ctx context.Context, userID uuid.UUID) (*entity.MaximumSpend, error)
}

type maximumSpendRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewMaximumSpendRepository(db *sql.DB, redis *redis.Client) MaximumSpendRepository {
	return &maximumSpendRepository{
		db:    db,
		redis: redis,
	}
}

func (r *maximumSpendRepository) CheckAlert(ctx context.Context, ms *entity.MaximumSpend, period string) error {
	query := `
		SELECT id, user_id, amount, period
		FROM maximum_spends
		WHERE user_id = $1 AND period = $2
	`
	var row *sql.Row
	if period == "daily" {
		row = r.db.QueryRowContext(ctx, query, ms.UserID, ms.DailyLimit)
	} else if period == "monthly" {
		row = r.db.QueryRowContext(ctx, query, ms.UserID, ms.MonthlyLimit)
	} else if period == "yearly" {
		row = r.db.QueryRowContext(ctx, query, ms.UserID, ms.YearlyLimit)
	}

	var existingMS entity.MaximumSpend
	err := row.Scan(&existingMS.ID, &existingMS.UserID, &existingMS.DailyLimit, &existingMS.MonthlyLimit, &existingMS.YearlyLimit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No existing maximum spend found
		}
		return err
	}

	ms.ID = existingMS.ID // Set the ID for update operation
	return nil
}

func (r *maximumSpendRepository) UpsertMaximumSpend(ctx context.Context, ms *entity.MaximumSpend) error {
	query := `
		INSERT INTO maximum_spends (id, user_id, daily_limit, monthly_limit, yearly_limit)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
		SET daily_limit = EXCLUDED.daily_limit,
			monthly_limit = EXCLUDED.monthly_limit,
			yearly_limit = EXCLUDED.yearly_limit,
			updated_at = NOW()
	`
	if ms.ID == "" {
		ms.ID = id.GenerateID()
	}

	_, err := r.db.ExecContext(ctx, query,
		ms.ID,
		ms.UserID,
		ms.DailyLimit,
		ms.MonthlyLimit,
		ms.YearlyLimit,
	)

	return err
}
func (r *maximumSpendRepository) GetMaximumSpendByUserID(ctx context.Context, userID uuid.UUID) (*entity.MaximumSpend, error) {
	query := `
		SELECT id, user_id, daily_limit, monthly_limit, yearly_limit
		FROM maximum_spends
		WHERE user_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, userID)
	ms := &entity.MaximumSpend{}
	err := row.Scan(
		&ms.ID,
		&ms.UserID,
		&ms.DailyLimit,
		&ms.MonthlyLimit,
		&ms.YearlyLimit,
	)

	if err == sql.ErrNoRows {
		return nil, errx.NewNotFoundError("User doesnt set maximum spend yet")
	}

	if err != nil {
		return nil, err
	}

	return ms, nil
}
