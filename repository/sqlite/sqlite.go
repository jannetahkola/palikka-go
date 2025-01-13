package sqlite

import (
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	sqlc "palikka-go/repository/sqlite/.sqlc"
	"slices"
)

//go:embed migration/*sql
var sqlMigrateFS embed.FS

type DB struct {
	db      *sql.DB
	dsn     string
	querier sqlc.Querier
}

func NewDB(dsn string) *DB {
	return &DB{
		dsn: dsn,
	}
}

func (db *DB) Open() (err error) {
	if db.db, err = sql.Open("sqlite", db.dsn); err != nil {
		return err
	}

	db.querier = sqlc.New(db.db)

	if err := db.migrate(); err != nil {
		return fmt.Errorf("sqlite: migrate: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	if err := db.db.Close(); err != nil {
		return err
	}
	return nil
}

func (db *DB) migrate() error {
	// Ensure the 'migrations' table exists so we don't duplicate migrations.
	if _, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("sqlite: create migrations table: %w", err)
	}

	names, err := fs.Glob(sqlMigrateFS, "migration/*.sql")
	if err != nil {
		return err
	}
	slices.Sort(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("sqlite: migrate file=%q: %w", name, err)
		}
		fmt.Printf("sqlite: migrate file=%q ok\n", name)
	}

	return nil
}

func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(sqlMigrateFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}
