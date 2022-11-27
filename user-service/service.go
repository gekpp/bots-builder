package user_service

import (
	"fmt"
	"github.com/google/uuid"
)

type (
	User struct {
		ID       uuid.UUID
		UserName string
	}

	UserService interface {
		CreateOrGetTelegramUser(userName string) (*User, error)
		GetTelegramUser(TelegramUserDescription) (*User, error)
	}

	TelegramUserDescription struct {
		ID uuid.UUID
		// ...
	}

	service struct {
		storage UserStorage
	}
)

func NewService(storage UserStorage) *service {
	return &service{
		storage: storage,
	}
}

func (s *service) CreateOrGetTelegramUser(userName string) (*User, error) {
	if userName == "" {
		return nil, fmt.Errorf("empty username")
	}

	strgUser, err := s.storage.Create(userName)
	if err != nil {
		return nil, fmt.Errorf("can not create user: %w", err)
	}

	return &User{
		ID:       strgUser.ID,
		UserName: strgUser.UserName,
	}, nil
}

func (s *service) GetTelegramUser(info TelegramUserDescription) (*User, error) {
	user, err := s.storage.GetByID(info.ID)
	if err != nil {
		return nil, fmt.Errorf("can not get user by id: %w", err)
	}

	return &User{
		ID:       user.ID,
		UserName: user.UserName,
	}, nil
}
