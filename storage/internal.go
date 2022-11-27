package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type (
	ConnectionInfo struct {
		Host     string
		Port     int
		Name     string
		Username string
		Password string
		Timeout  int
		SSLMode  string
	}
)

func Connect(info ConnectionInfo) *sql.DB {
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		info.Host, info.Port, info.Username, info.Password, info.Name, info.SSLMode, info.Timeout,
	)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(0)

	if err := db.Ping(); err != nil {
		log.Fatalf("database is not available right now: %v", err)
	}

	return db
}
