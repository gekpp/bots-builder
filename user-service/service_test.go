package users

import (
	st "github.com/gekpp/bots-builder/internal/infra"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

// TestService Was made for local tests
func TestService(t *testing.T) {
	t.Skip()
	con := st.ConnectDB(
		"localhost",
		5432,
		"tgbot",
		"postgres",
		"postgres",
		15,
		"disable",
	)
	storage, err := NewStorage(con)
	require.NoError(t, err)

	telegramID := strconv.FormatInt(time.Now().Unix(), 10)
	service := NewService(storage)

	id := uuid.New()

	description := TelegramUserDescription{
		ID:         id,
		TelegramID: telegramID,
		FirstName:  "test",
		LastName:   "test",
		UserName:   "test",
	}

	expectedUser := User{
		ID:         id,
		TelegramID: telegramID,
		FirstName:  "test",
		LastName:   "test",
		UserName:   "test",
	}

	// first try
	user, err := service.CreateOrGetTelegramUser(description)
	require.NoError(t, err)
	require.Equal(t, expectedUser, *user)

	//try to create existing user
	user, err = service.CreateOrGetTelegramUser(description)
	require.NoError(t, err)
	require.Equal(t, expectedUser, *user)

	// get user
	user, err = service.GetUserByID(description)
	require.NoError(t, err)
	require.Equal(t, expectedUser, *user)

	user, err = service.GetUserByTelegramID(description)
	require.NoError(t, err)
	require.Equal(t, expectedUser, *user)

	_, err = service.GetUserByID(TelegramUserDescription{ID: uuid.New()})
	require.ErrorIs(t, err, ErrNotFound)
}
