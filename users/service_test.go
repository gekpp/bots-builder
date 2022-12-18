package users

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gekpp/bots-builder/internal/tests"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var shutdownFunc func()

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	// order matters
	binds := []string{
		fmt.Sprintf("%s/../scripts/migrations/20221210-init-schema.sql", workingDir),         // 1
		fmt.Sprintf("%s/../scripts/migrations/data/20221212-test-init-data.sql", workingDir), // 1
	}

	bindsMap := map[string]string{}

	names := []rune{'a', 'a', 'a'}
	k := len(names) - 1
	for _, s := range binds {
		bindsMap[s] = fmt.Sprintf("/docker-entrypoint-initdb.d/%s.sql", string(names))
		for i := k; i <= k; i-- {
			if names[i] < 'z' { // if panics increase names slice
				names[i]++
				break
			}
			if names[i] == 'z' {
				names[i] = 'a'
			}
		}
	}
	db, shutdownFunc = tests.ConnectTestContainers(bindsMap)
	defer shutdownFunc()

	os.Exit(m.Run())
}

func TestService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	telegramID := strconv.FormatInt(time.Now().Unix(), 10)
	service := New(sqlx.NewDb(db, "postgres"))

	id := uuid.New()

	tgUser := User{
		ID:               id,
		TelegramID:       telegramID,
		FirstName:        "test",
		LastName:         "test",
		TelegramUserName: "test",
	}

	expectedUser := User{
		ID:               id,
		TelegramID:       telegramID,
		FirstName:        "test",
		LastName:         "test",
		TelegramUserName: "test",
	}

	// first try
	expectedUser, err := service.CreateOrGetTelegramUser(ctx, tgUser)
	require.NoError(t, err)
	// require.Equal(t, expectedUser, user)

	//try to create existing user
	user, err := service.CreateOrGetTelegramUser(ctx, tgUser)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// get user
	user, err = service.GetByID(ctx, expectedUser.ID)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	user, err = service.GetByTelegramID(ctx, tgUser.TelegramID)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	_, err = service.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, ErrNotFound)
}
