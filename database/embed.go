package database

import (
	_ "embed"
)

//go:embed schema.sql
var schema string

func Schema() string {
	return schema
}
