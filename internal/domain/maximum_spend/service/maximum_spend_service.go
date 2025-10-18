package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kenziehh/cashflow-be/internal/domain/maximum_spend/entity"
	"github.com/kenziehh/cashflow-be/internal/domain/maximum_spend/repository"
	"github.com/kenziehh/cashflow-be/pkg/errx"
)

type MaximumSpendService interface {
	SetMaximumSpend(ctx context.Context, userID uuid.UUID, daily, monthly, yearly float64) (*entity.MaximumSpend, error)
	GetMaximumSpend(ctx context.Context, userID uuid.UUID) (*entity.MaximumSpend, error)
}

type maximumSpendService struct {
	repo repository.MaximumSpendRepository
}

func NewMaximumSpendService(repo repository.MaximumSpendRepository) MaximumSpendService {
	return &maximumSpendService{
		repo: repo,
	}
}

func (s *maximumSpendService) SetMaximumSpend(ctx context.Context, userID uuid.UUID, daily, monthly, yearly float64) (*entity.MaximumSpend, error) {
	// Cek apakah user sudah pernah set sebelumnya
	existing, err := s.repo.GetMaximumSpendByUserID(ctx, userID)
	if err != nil {
		if err == errx.NewNotFoundError("User doesnt set maximum spend yet") {
			// Belum ada, buat baru
			newMS := &entity.MaximumSpend{
				ID:           uuid.NewString(),
				UserID:       userID,
				DailyLimit:   daily,
				MonthlyLimit: monthly,
				YearlyLimit:  yearly,
			}
			if err := s.repo.UpsertMaximumSpend(ctx, newMS); err != nil {
				return nil, err
			}
			return newMS, nil
		}
		return nil, err
	}

	// Sudah ada, update data
	existing.DailyLimit = daily
	existing.MonthlyLimit = monthly
	existing.YearlyLimit = yearly

	if err := s.repo.UpsertMaximumSpend(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *maximumSpendService) GetMaximumSpend(ctx context.Context, userID uuid.UUID) (*entity.MaximumSpend, error) {
	ms, err := s.repo.GetMaximumSpendByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
