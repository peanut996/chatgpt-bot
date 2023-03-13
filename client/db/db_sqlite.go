package db

import (
	"chatgpt-bot/cfg"
	"database/sql"
)

type SQLiteDB struct {
	db         *sql.DB
	dbLockChan chan interface{}
}

func NewSQLiteDB() *SQLiteDB {
	return &SQLiteDB{}
}

func (s *SQLiteDB) Init(cfg *cfg.Config) error {
	var dbPath string
	if cfg.DatabaseConfig.Path != "" {
		dbPath = cfg.DatabaseConfig.Path
	} else {
		dbPath = "./bot.db"
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	s.db = db

	s.dbLockChan = make(chan interface{}, 1)
	return nil
}

func (s *SQLiteDB) Lock() {
	s.dbLockChan <- struct{}{}
}

func (s *SQLiteDB) Unlock() {
	<-s.dbLockChan
}

func (s *SQLiteDB) Close() {
	s.db.Close()
}

func (s *SQLiteDB) Query(query string, args ...any) (*sql.Rows, error) {
	s.Lock()
	defer s.Unlock()
	return s.db.Query(query, args...)
}

func (s *SQLiteDB) QueryRow(query string, args ...any) *sql.Row {
	s.Lock()
	defer s.Unlock()
	return s.db.QueryRow(query, args...)
}

func (s *SQLiteDB) Exec(query string, args ...any) (sql.Result, error) {
	s.Lock()
	defer s.Unlock()
	return s.db.Exec(query, args...)
}
