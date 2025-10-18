package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/kenziehh/cashflow-be/internal/domain/auth/dto"
	"github.com/kenziehh/cashflow-be/internal/domain/auth/entity"
	"github.com/kenziehh/cashflow-be/pkg/errx"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	StoreToken(ctx context.Context, userID uuid.UUID, token string, expiration time.Duration) error
	DeleteToken(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error
}

type authRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewAuthRepository(db *sql.DB, redis *redis.Client) AuthRepository {
	return &authRepository{
		db:    db,
		redis: redis,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return errx.ErrDatabaseError
	}

	return nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errx.ErrUserNotFound
	}

	if err != nil {
		return nil, errx.ErrDatabaseError
	}

	return user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errx.ErrUserNotFound
	}

	if err != nil {
		return nil, errx.ErrDatabaseError
	}

	return user, nil
}

func (r *authRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error {
	query := `
		UPDATE users
		SET name = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, req.Name, time.Now(), userID)
	if err != nil {
		return errx.ErrDatabaseError
	}

	return nil
}

func (r *authRepository) StoreToken(ctx context.Context, userID uuid.UUID, token string, expiration time.Duration) error {
	key := "token:" + token
	err := r.redis.Set(ctx, key, userID.String(), expiration).Err()
	if err != nil {
		return errx.ErrRedisError
	}
	return nil
}

func (r *authRepository) DeleteToken(ctx context.Context, token string) error {
	key := "blacklist:" + token
	err := r.redis.Set(ctx, key, "1", 24*time.Hour).Err()
	if err != nil {
		return errx.ErrRedisError
	}
	return nil
}

func (r *authRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := "blacklist:" + token
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, errx.ErrRedisError
	}
	return val == "1", nil
}
