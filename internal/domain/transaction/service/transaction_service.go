package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kenziehh/cashflow-be/internal/domain/transaction/dto"
	"github.com/kenziehh/cashflow-be/internal/domain/transaction/entity"
	"github.com/kenziehh/cashflow-be/internal/domain/transaction/repository"
	"github.com/kenziehh/cashflow-be/pkg/errx"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID uuid.UUID) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, req dto.UpdateTransactionRequest) (*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetTransactionsWithPagination(ctx context.Context, userID uuid.UUID, params dto.TransactionListParams) (dto.PaginatedTransactionsResponse, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID uuid.UUID) (*entity.Transaction, error) {
	tx := &entity.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		CategoryID:      req.CategoryID,
		Note:            req.Note,
		Date:            req.Date,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	tx, err := s.repo.GetTransactionByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errx.ErrTransactionNotFound
	}
	return tx, nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, id uuid.UUID, req dto.UpdateTransactionRequest) (*entity.Transaction, error) {
	tx, err := s.repo.GetTransactionByID(ctx, id.String())
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errx.ErrTransactionNotFound
	}

	// Update fields
	if req.Amount != 0 {
		tx.Amount = req.Amount
	}
	if req.TransactionType != "" {
		tx.TransactionType = req.TransactionType
	}
	if req.CategoryID != "" {
		tx.CategoryID = req.CategoryID
	}
	if req.Note != "" {
		tx.Note = req.Note
	}
	if req.Date != "" {
		tx.Date = req.Date
	}
	tx.UpdatedAt = time.Now()

	if err := s.repo.UpdateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	tx, err := s.repo.GetTransactionByID(ctx, id.String())
	if err != nil {
		return err
	}
	if tx == nil {
		return errx.ErrTransactionNotFound
	}

	if err := s.repo.DeleteTransaction(ctx, id.String()); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) GetTransactionsWithPagination(ctx context.Context, userID uuid.UUID, params dto.TransactionListParams) (dto.PaginatedTransactionsResponse, error) {
	txs, err := s.repo.GetTransactionsWithPagination(ctx, userID, params)
	if err != nil {
		return dto.PaginatedTransactionsResponse{}, err
	}
	return txs, nil
}						

