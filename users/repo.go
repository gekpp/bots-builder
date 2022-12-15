package users

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	repository struct {
		db *sqlx.DB
	}

	user struct {
		ID               uuid.UUID `db:"id"`
		TelegramID       string    `db:"telegram_user_id"`
		FirstName        string    `db:"first_name"`
		LastName         string    `db:"last_name"`
		TelegramUserName string    `db:"telegram_user_name"`
	}
)

func newRepository(db *sql.DB) *repository {
	r := repository{
		db: sqlx.NewDb(db, "postgres").Unsafe(),
	}

	return &r
}

func (s *repository) getOrCreate(ctx context.Context, tgUser User) (user, error) {
	res := user{}
	query := `INSERT INTO users (id,  telegram_user_id, telegram_user_name, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_user_id) DO UPDATE
		SET telegram_user_name = excluded.telegram_user_name,
			first_name = excluded.first_name,
			last_name = excluded.last_name
		RETURNING id, telegram_user_id, telegram_user_name, first_name, last_name;
	`
	err := s.db.GetContext(ctx, &res, query, tgUser.ID.String(), tgUser.TelegramID, tgUser.TelegramUserName, tgUser.FirstName, tgUser.LastName)
	if err == sql.ErrNoRows {
		return user{}, ErrNotFound
	}
	return res, err
}

func (s *repository) getByID(ctx context.Context, id uuid.UUID) (user, error) {
	res := user{}
	err := s.db.GetContext(ctx, &res,
		`SELECT id, telegram_user_id, telegram_user_name, first_name, last_name 
FROM users WHERE id = $1`, id.String())
	if err == sql.ErrNoRows {
		return user{}, ErrNotFound
	}
	return res, err
}

func (s *repository) getByTelegramID(ctx context.Context, telegramID string) (user, error) {
	res := user{}
	err := s.db.GetContext(ctx, &res, `SELECT id, telegram_user_id, telegram_user_name, first_name, last_name 
FROM users WHERE telegram_user_id = $1`, telegramID)
	if err == sql.ErrNoRows {
		return user{}, ErrNotFound
	}
	return res, err
}
