package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

const (
	ConnectionString = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
)

type Database struct {
	conn *sql.DB
}

var (
	instance *Database
	once     sync.Once
	initErr  error
)

// GetDB returns the singleton database instance, initializing it on first call
func GetDB() (*sql.DB, error) {
	once.Do(func() {
		slog.Info("Initializing database connection")

		dbUser := envOrDefault("DB_USER", "postgres")
		dbPass := envOrDefault("DB_PASS", "password")
		dbHost := envOrDefault("DB_HOST", "localhost")
		dbPort := envOrDefault("DB_PORT", "5432")
		dbName := envOrDefault("DB_NAME", "database")

		conString := fmt.Sprintf(ConnectionString, dbUser, dbPass, dbHost, dbPort, dbName)

		conn, err := sql.Open("postgres", conString)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err = conn.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		instance = &Database{conn: conn}
		slog.Info("Database connection established successfully")
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance.conn, nil
}

func envOrDefault(env, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}
	return value
}

// Close closes the database connection
func Close() error {
	if instance != nil && instance.conn != nil {
		slog.Info("Closing database connection")
		return instance.conn.Close()
	}
	return nil
}
