package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type (
	User struct {
		ID               uuid.UUID
		TelegramID       string
		FirstName        string
		LastName         string
		TelegramUserName string
	}
)

var (
	ErrNotFound = errors.New("not found")
)

type Service interface {
	CreateOrGetTelegramUser(ctx context.Context, info User) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByTelegramID(ctx context.Context, telegramID string) (User, error)
}
