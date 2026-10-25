package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func openTxTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.ExecContext(context.Background(), "CREATE TABLE tx_values (v TEXT NOT NULL)"); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return db
}

func TestWithTx_CommitsOnSuccess(t *testing.T) {
	db := openTxTestDB(t)
	ctx := context.Background()

	err := WithTx(ctx, db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO tx_values(v) VALUES (?)", "ok")
		return err
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tx_values WHERE v = ?", "ok").Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected committed row, got %d", count)
	}
}

func TestWithTx_RollsBackAndWrapsError(t *testing.T) {
	db := openTxTestDB(t)
	ctx := context.Background()
	fnErr := errors.New("boom")

	err := WithTx(ctx, db, func(tx *sqlx.Tx) error {
		if _, e := tx.ExecContext(ctx, "INSERT INTO tx_values(v) VALUES (?)", "rollback"); e != nil {
			return e
		}
		return fnErr
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "tx: boom") {
		t.Fatalf("expected wrapped tx error, got %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tx_values WHERE v = ?", "rollback").Scan(&count); err != nil && !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("count rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected rollback, found %d rows", count)
	}
}
