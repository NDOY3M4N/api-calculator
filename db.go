package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabaseSqlite(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}
