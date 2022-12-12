package users

import (
	"fmt"
	"github.com/google/uuid"
)

type (
	User struct {
		ID         uuid.UUID
		TelegramID string
		FirstName  string
		LastName   string
		UserName   string
	}

	TelegramUserDescription struct {
		ID         uuid.UUID
		TelegramID string
		FirstName  string
		LastName   string
		UserName   string
	}

	service struct {
		repo UserStorage
	}
)

func NewService(storage UserStorage) *service {
	return &service{
		repo: storage,
	}
}

func (s *service) CreateOrGetTelegramUser(info TelegramUserDescription) (*User, error) {
	if info.TelegramID == "" {
		return nil, fmt.Errorf("empty telegram id")
	}

	strgUser, err := s.repo.GetOrCreate(info)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:         strgUser.ID,
		TelegramID: strgUser.TelegramID,
		FirstName:  strgUser.FirstName,
		LastName:   strgUser.LastName,
		UserName:   strgUser.UserName,
	}, nil
}

func (s *service) GetUserByID(info TelegramUserDescription) (*User, error) {
	user, err := s.repo.GetByID(info.ID)
	if err != nil {
		return nil, fmt.Errorf("can not get user by id: %w", err)
	}

	return &User{
		ID:         user.ID,
		TelegramID: user.TelegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserName:   user.UserName,
	}, nil
}

func (s *service) GetUserByTelegramID(info TelegramUserDescription) (*User, error) {
	user, err := s.repo.GetByTelegramID(info.TelegramID)
	if err != nil {
		return nil, fmt.Errorf("can not get user by telegram_id: %w", err)
	}

	return &User{
		ID:         user.ID,
		TelegramID: user.TelegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserName:   user.UserName,
	}, nil
}
