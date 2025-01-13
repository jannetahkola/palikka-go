package config

import _ "embed"

//go:embed data/schema.sql
var createSchemaSql string

//go:embed data/reset.sql
var resetSchemaSql string

//go:embed data/seed.sql
var seedSql string

func CreateSchemaSql() string {
	return createSchemaSql
}

func ResetSchemaSql() string {
	return resetSchemaSql
}

func SeedSql() string {
	return seedSql
}
