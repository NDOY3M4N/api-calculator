package main

import (
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewDatabaseSqlite(name string) (*sql.DB, error) {
	if _, err := os.Stat(name); errors.Is(err, fs.ErrNotExist) {
		file, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		file.Close()

		db, err := openConnection(name)
		if err != nil {
			return nil, err
		}

		goose.SetDialect("sqlite3")
		goose.SetBaseFS(embedMigrations)

		if err := goose.Up(db, "migrations"); err != nil {
			return nil, err
		}

		return db, nil
	} else if err != nil {
		return nil, err
	}

	return openConnection(name)
}

func openConnection(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	return db, nil
}
