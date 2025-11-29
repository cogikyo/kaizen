package db

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

var conn *sqlx.DB

var ErrNotKaizen = errors.New("not a kaizen directory")

func Dir() string {
	return ".kaizen"
}

func Path() string {
	return filepath.Join(Dir(), "kaizen.db")
}

func Exists() bool {
	_, err := os.Stat(Dir())
	return err == nil
}

func Open() error {
	if conn != nil {
		return nil
	}

	if !Exists() {
		return ErrNotKaizen
	}

	var err error
	conn, err = sqlx.Open("sqlite", Path())
	if err != nil {
		return err
	}

	_, err = conn.Exec(schema)
	if err != nil {
		conn.Close()
		conn = nil
		return err
	}

	return nil
}

func Create() error {
	if Exists() {
		return nil
	}

	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return err
	}

	return Open()
}

func Get() *sqlx.DB {
	return conn
}

func Close() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}
