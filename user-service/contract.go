package users

import "github.com/google/uuid"

type (
	UserStorage interface {
		GetOrCreate(user TelegramUserDescription) (*StorageUser, error)
		GetByID(id uuid.UUID) (*StorageUser, error)
		GetByTelegramID(telegramID string) (*StorageUser, error)
	}
)
