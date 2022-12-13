package users

import (
	"context"
	st "github.com/gekpp/bots-builder/internal/infra"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

// TestService Was made for local tests
func TestService(t *testing.T) {
	//t.Skip()
	ctx := context.Background()
	con := st.ConnectDB(
		"localhost",
		5432,
		"tgbot",
		"postgres",
		"postgres",
		15,
		"disable",
	)

	telegramID := strconv.FormatInt(time.Now().Unix(), 10)
	service, err := NewService(con)
	require.NoError(t, err)

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
	user, err := service.CreateOrGetTelegramUser(ctx, description)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	//try to create existing user
	user, err = service.CreateOrGetTelegramUser(ctx, description)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	// get user
	user, err = service.GetByID(ctx, description.ID)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	user, err = service.GetByTelegramID(ctx, description.TelegramID)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	_, err = service.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, ErrNotFound)
}
