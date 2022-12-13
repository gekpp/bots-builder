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

func NewService(db *sql.DB) (*service, error) {
	repo, err := newRepository(db)
	if err != nil {
		return nil, err
	}

	return &service{
		repo: repo,
	}, nil
}

func (s *service) CreateOrGetTelegramUser(ctx context.Context, info TelegramUserDescription) (User, error) {
	if info.TelegramID == "" {
		return User{}, fmt.Errorf("empty telegram id")
	}

	strgUser, err := s.repo.getOrCreate(info)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:         strgUser.ID,
		TelegramID: strgUser.TelegramID,
		FirstName:  strgUser.FirstName,
		LastName:   strgUser.LastName,
		UserName:   strgUser.UserName,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := s.repo.getByID(id)
	if err != nil {
		return User{}, fmt.Errorf("can not get user by id: %w", err)
	}

	return User{
		ID:         user.ID,
		TelegramID: user.TelegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserName:   user.UserName,
	}, nil
}

func (s *service) GetByTelegramID(ctx context.Context, telegramID string) (User, error) {
	user, err := s.repo.getByTelegramID(telegramID)
	if err != nil {
		return User{}, fmt.Errorf("can not get user by telegram_id: %w", err)
	}

	return User{
		ID:         user.ID,
		TelegramID: user.TelegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserName:   user.UserName,
	}, nil
}
