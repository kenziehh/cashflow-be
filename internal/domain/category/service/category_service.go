package service

import (
	

	"context"

	"github.com/kenziehh/cashflow-be/internal/domain/category/repository"
	"github.com/kenziehh/cashflow-be/internal/domain/category/dto"
)

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]dto.GetAllCategoryResponse, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]dto.GetAllCategoryResponse, error) {
	categories, err := s.repo.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}