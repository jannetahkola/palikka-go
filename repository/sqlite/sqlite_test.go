package sqlite_test

import (
	"palikka-go/repository/sqlite"
	"testing"
)

func TestDB(t *testing.T) {
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

func MustOpenDB(t *testing.T) *sqlite.DB {
	t.Helper()

	db := sqlite.NewDB(":memory:")
	if err := db.Open(); err != nil {
		t.Fatal(err)
	}

	return db
}

func MustCloseDB(t *testing.T, db *sqlite.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
