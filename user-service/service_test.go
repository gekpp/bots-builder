package users

import (
	"context"
	"database/sql"
	"fmt"
	st "github.com/gekpp/bots-builder/internal/infra"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"testing"
	"time"
)

func NewTestDatabase(ctx context.Context) (con *sql.DB, terminateFunc func(context.Context) error, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	terminateFunc = postgres.Terminate
	host, err := postgres.Host(ctx)
	if err != nil {
		return
	}

	natPort, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		return
	}

	p, _ := strconv.Atoi(natPort.Port())
	con = st.ConnectDB(
		host,
		p,
		"postgres",
		"postgres",
		"postgres",
		15,
		"disable",
	)

	return
}

func CreateUsersTable(db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE %s (
		id uuid NOT NULL primary key,
		telegram_id varchar(50) unique,
		first_name varchar(50),
		last_name varchar(50),
		user_name varchar(50)
	);`, usersTableName)
	_, err := db.Exec(query)
	return err
}

func TestService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	con, terminateFunc, err := NewTestDatabase(ctx)
	defer terminateFunc(ctx)
	require.NoError(t, err)

	err = CreateUsersTable(con)
	require.NoError(t, err)

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
