package users

import (
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
