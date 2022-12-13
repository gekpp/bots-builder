package users

import (
	"errors"
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
)

var (
	ErrNotFound = errors.New("not found")
)
