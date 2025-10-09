package repository

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/kenziehh/cashflow-be/internal/domain/category/dto"
	"github.com/kenziehh/cashflow-be/pkg/errx"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]dto.GetAllCategoryResponse, error)
}

type categoryRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewCategoryRepository(db *sql.DB, redis *redis.Client) CategoryRepository {
	return &categoryRepository{
		db:    db,
		redis: redis,
	}
}

func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]dto.GetAllCategoryResponse, error) {
	query := `
	SELECT id, name
	FROM categories
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errx.ErrDatabaseError
	}
	defer rows.Close()

	var categories []dto.GetAllCategoryResponse

	for rows.Next() {
		var c dto.GetAllCategoryResponse
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, errx.ErrDatabaseError
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, errx.ErrDatabaseError
	}

	return categories, nil
}
