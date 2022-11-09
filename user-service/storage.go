package user_service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type (
	storage struct {
		db *sql.DB
	}

	User struct {
		ID       uuid.UUID
		UserName string
	}

	UserStorage interface {
		Create(userName string) (*User, error)
		GetByName(username string) (*User, error)
		GetByID(id uuid.UUID) (*User, error)
	}
)

const (
	userTableName = "users"
)

var (
	ErrNotFound = errors.New("not_found")
)

func NewStorage(db *sql.DB) (*storage, error) {
	s := storage{
		db: db,
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to DB: %v", err)
	}

	if _, err := db.Exec(fmt.Sprintf("SELECT 1 FROM %q LIMIT 1", userTableName)); err != nil {
		return nil, fmt.Errorf("could not connect to storage: %v", err)
	}

	return &s, nil
}

func (s *storage) Create(userName string) (*User, error) {
	user, err := s.GetByName(userName)
	if err != nil {
		return user, nil
	}
	if err != ErrNotFound {
		return nil, fmt.Errorf("error on checking user before creating: %v", err)
	}

	u := &User{
		ID:       uuid.New(),
		UserName: userName,
	}
	if err = s.create(u); err != nil {
		return nil, fmt.Errorf("can not create user: %v", err)
	}

	return u, nil
}

func (s *storage) GetByName(username string) (*User, error) {
	queryRow := s.db.QueryRow(
		fmt.Sprintf(`SELECT id, user_name FROM %s WHERE user_name = $1`, userTableName), username)

	return s.getFromQueryRow(queryRow)
}

func (s *storage) GetByID(id uuid.UUID) (*User, error) {
	queryRow := s.db.QueryRow(
		fmt.Sprintf(`SELECT id, user_name FROM %s WHERE id = $1`, userTableName), id.String())

	return s.getFromQueryRow(queryRow)
}

func (s *storage) getFromQueryRow(query *sql.Row) (*User, error) {
	user := &User{}
	err := query.Scan(user.ID, user.UserName)
	if err == nil {
		return user, nil
	}
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return nil, fmt.Errorf("can not get user: %v", err)
}

func (s *storage) create(user *User) error {
	_, err := s.db.Exec(
		fmt.Sprintf(`INSERT INTO %s (id, user_name) VALUES ($1, $2)`, userTableName), user.ID.String(), user.UserName)
	return err
}
