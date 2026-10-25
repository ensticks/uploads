// Package database provides PostgreSQL connectivity via sqlx.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // register "pgx" driver
	"github.com/jmoiron/sqlx"
)

// Config holds the settings required to open a database connection.
type Config struct {
	// DSN is the PostgreSQL connection string.
	DSN string
}

// Open establishes a sqlx PostgreSQL connection using the settings in cfg.
func Open(cfg Config, log *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("database: open: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	log.Info("database connected")
	return db, nil
}
