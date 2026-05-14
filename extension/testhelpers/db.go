package testhelpers

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func SetupTestDB(t *testing.T) *sql.Tx {
	dsn := os.Getenv("DATABASE_URL")
	// TODO: should move to a config file
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/post_confiavel?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	return tx
}
