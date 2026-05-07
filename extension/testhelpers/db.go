package testhelpers

import (
	"database/sql"
	"testing"
)

func SetupTestDB(t *testing.T) *sql.Tx {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5433/post_con_back?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	return tx
}