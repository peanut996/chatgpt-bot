package db

import (
	"chatgpt-bot/cfg"
	"database/sql"
)

type BotDB interface {
	Init(*cfg.Config) error
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}
