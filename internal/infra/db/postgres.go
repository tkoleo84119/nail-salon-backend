package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
)

type Database struct {
	Std  *sql.DB
	Sqlx *sqlx.DB
}

func New(dbConfig config.DBConfig) (*Database, func(), error) {
	db, err := sql.Open("postgres", dbConfig.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(dbConfig.ConnMaxLife)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("ping db: %w", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	cleanup := func() { _ = db.Close() }

	return &Database{
		Std:  db,
		Sqlx: sqlxDB,
	}, cleanup, nil
}
