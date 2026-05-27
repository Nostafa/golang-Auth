package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func New(dbUrl string, maxOpenConns, maxIdleConns int, maxLifetime, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	lifetime, err := time.ParseDuration(maxLifetime)
	if err != nil {
		return nil, err
	}
	idleTime, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(lifetime)
	db.SetConnMaxIdleTime(idleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil

}
