package users

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	repo struct {
		db *sqlx.DB
	}

	StorageUser struct {
		ID         uuid.UUID `db:"id"`
		TelegramID string    `db:"telegram_id"`
		FirstName  string    `db:"first_name"`
		LastName   string    `db:"last_name"`
		UserName   string    `db:"user_name"`
	}
)

const (
	usersTableName = "users"
)

var (
	ErrNotFound = errors.New("not found")
)

func NewStorage(db *sql.DB) (*repo, error) {
	s := repo{
		db: sqlx.NewDb(db, usersTableName),
	}

	if err := s.db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to DB: %w", err)
	}

	if _, err := s.db.Exec(fmt.Sprintf("SELECT 1 FROM %q LIMIT 1", s.db.DriverName())); err != nil {
		return nil, fmt.Errorf("could not check %s table exists: %w", s.db.DriverName(), err)
	}

	return &s, nil
}

func (s *repo) GetOrCreate(user TelegramUserDescription) (*StorageUser, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (id,  telegram_id, user_name, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_id) DO UPDATE
		SET user_name = excluded.user_name,
			first_name = excluded.first_name,
			last_name = excluded.last_name
		RETURNING id, telegram_id, user_name, first_name, last_name;
	`, s.db.DriverName())
	queryRow := s.db.QueryRowx(
		query, user.ID.String(), user.TelegramID, user.UserName, user.FirstName, user.LastName)

	s.db.Driver()
	return getFromQueryRow(queryRow)
}

func (s *repo) GetByID(id uuid.UUID) (*StorageUser, error) {
	queryRow := s.db.QueryRowx(
		fmt.Sprintf(`SELECT id, telegram_id, user_name, first_name, last_name FROM %s WHERE id = $1`, s.db.DriverName()), id.String())

	return getFromQueryRow(queryRow)
}

func (s *repo) GetByTelegramID(TelegramID string) (*StorageUser, error) {
	queryRow := s.db.QueryRowx(
		fmt.Sprintf(`SELECT id, telegram_id, user_name, first_name, last_name FROM %s WHERE telegram_id = $1`, s.db.DriverName()), TelegramID)

	return getFromQueryRow(queryRow)
}

func getFromQueryRow(query *sqlx.Row) (*StorageUser, error) {
	user := &StorageUser{}
	err := query.StructScan(user)
	if err == nil {
		return user, nil
	}
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return nil, err
}
