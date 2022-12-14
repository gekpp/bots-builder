package infra

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type (
	Option func(*sql.DB)
)

func ConnectDB(
	Host string,
	Port int,
	Name string,
	Username string,
	Password string,
	Timeout int,
	SSLMode string,
	opts ...Option,
) *sql.DB {
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		Host, Port, Username, Password, Name, SSLMode, Timeout,
	)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("could not connect to database: %w", err)
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
