package test

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"palikka-go/config"
	sqlc2 "palikka-go/repository/sqlite/.sqlc"
)

// todo exclude from build?

func InitDatabase() (*sql.DB, error) {
	fmt.Println("database init")

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ResetDatabase(db *sql.DB) error {
	fmt.Println("database reset")

	_, err := db.Exec(config.ResetSchemaSql())
	if err != nil {
		return err
	}
	_, err = db.Exec(config.CreateSchemaSql())
	if err != nil {
		return err
	}
	return nil
}

func SeedDatabase(db *sql.DB) error {
	fmt.Println("database seed")

	_, err := db.Exec(config.SeedSql())
	if err != nil {
		return err
	}
	return nil
}

func InsertRolesTx(
	ctx context.Context,
	db *sql.DB,
	queries *sqlc2.Queries,
	args ...sqlc2.InsertRoleParams,
) ([]int64, error) {

	var ids []int64
	tx, err := db.Begin()
	if err != nil {
		return ids, err
	}
	defer tx.Rollback()
	for _, arg := range args {
		id, err := queries.WithTx(tx).InsertRole(ctx, arg)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	err = tx.Commit()

	return ids, err
}
