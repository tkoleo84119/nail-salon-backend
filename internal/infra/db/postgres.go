package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
)

type Database struct {
	Std     *sql.DB
	Sqlx    *sqlx.DB
	PgxPool *pgxpool.Pool
}

func New(dbCfg config.DBConfig) (*Database, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(dbCfg.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("parse DSN: %w", err)
	}

	// set pgxpool config
	poolCfg.MaxConns = int32(dbCfg.MaxOpenConns)
	poolCfg.MaxConnLifetime = dbCfg.MaxConnMaxLife
	poolCfg.MaxConnLifetimeJitter = dbCfg.MaxConnLifetimeJitter // avoid same time reconnect
	poolCfg.MaxConnIdleTime = dbCfg.MaxConnIdleTime

	pgxPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("new pgx pool: %w", err)
	}
	if err := pgxPool.Ping(ctx); err != nil {
		pgxPool.Close()
		return nil, nil, fmt.Errorf("ping pgx pool: %w", err)
	}

	stdDB := stdlib.OpenDBFromPool(pgxPool)
	sqlxDB := sqlx.NewDb(stdDB, "pgx")

	cleanup := func() {
		_ = stdDB.Close()
		pgxPool.Close()
	}

	return &Database{
		Std:     stdDB,
		Sqlx:    sqlxDB,
		PgxPool: pgxPool,
	}, cleanup, nil
}
