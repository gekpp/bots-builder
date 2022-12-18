package infra

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type (
	Option func(*sql.DB)
)

// MustConnectDB connects to DB
func MustConnectDB(
	host string,
	port int,
	name string,
	username string,
	password string,
	timeout int,
	sslMode string,
	opts ...Option,
) *sql.DB {
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		host, port, username, password, name, sslMode, timeout,
	)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(0)

	for _, opt := range opts {
		opt(db)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("database is not available right now: %v", err)
	}

	return db
}

func WithMaxOpenConns(n int) Option {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(n)
	}
}

func WithMaxIdleConns(n int) Option {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(n)
	}
}
