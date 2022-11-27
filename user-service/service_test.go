package user_service

import (
	st "github.com/gekpp/bots-builder/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

// TestService Was made for local tests
func TestService(t *testing.T) {
	t.Skip()
	con := st.Connect(st.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Name:     "tgbot",
		Username: "postgres",
		Password: "postgres",
		Timeout:  15,
		SSLMode:  "disable",
	})
	storage, err := NewStorage(con)
	require.NoError(t, err)

	name := strconv.FormatInt(time.Now().Unix(), 10)
	service := NewService(storage)

	user, err := service.CreateOrGetTelegramUser(name)
	require.NoError(t, err)

	user, err = service.GetTelegramUser(TelegramUserDescription{ID: user.ID})
	require.NoError(t, err)

	require.Equal(t, name, user.UserName)

	_, err = service.GetTelegramUser(TelegramUserDescription{ID: uuid.New()})
	require.ErrorIs(t, err, ErrNotFound)
}
