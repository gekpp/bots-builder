package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type (
	service struct {
		repo *repository
	}
)

// New creates new service
func New(db *sql.DB) *service {
	repo := newRepository(db)

	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrGetTelegramUser(ctx context.Context, info User) (User, error) {
	if info.TelegramID == "" {
		return User{}, fmt.Errorf("empty telegram id")
	}

	u, err := s.repo.getOrCreate(ctx, info)
	if err != nil {
		return User{}, err
	}

	return User(u), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	u, err := s.repo.getByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("can not get user by id: %w", err)
	}

	return User(u), nil
}

func (s *service) GetByTelegramID(ctx context.Context, telegramID string) (User, error) {
	u, err := s.repo.getByTelegramID(ctx, telegramID)
	if err != nil {
		return User{}, fmt.Errorf("can not get user by telegram_id: %w", err)
	}

	return User(u), nil
}
