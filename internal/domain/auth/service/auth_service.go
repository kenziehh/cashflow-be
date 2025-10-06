package service

import (
	"context"
	"time"

	"github.com/kenziehh/cashflow-be/internal/domain/auth/dto"
	"github.com/kenziehh/cashflow-be/internal/domain/auth/entity"
	"github.com/kenziehh/cashflow-be/internal/domain/auth/repository"
	"github.com/kenziehh/cashflow-be/pkg/errx"

	"github.com/kenziehh/cashflow-be/pkg/bcrypt"
	"github.com/kenziehh/cashflow-be/pkg/jwt"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, token string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := s.repo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errx.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	// Create user
	user := &entity.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	// Store token in Redis
	if err := s.repo.StoreToken(ctx, user.ID, token, 24*time.Hour); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserProfile{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errx.ErrInvalidCredentials
	}

	// Verify password
	if !bcrypt.CheckPassword(req.Password, user.Password) {
		return nil, errx.ErrInvalidCredentials
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	// Store token in Redis
	if err := s.repo.StoreToken(ctx, user.ID, token, 24*time.Hour); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserProfile{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteToken(ctx, token)
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfile{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
