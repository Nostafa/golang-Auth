package main

import (
	"log"

	"github.com/Nostafa/golang-jwt/internal/database"
	"github.com/Nostafa/golang-jwt/internal/env"
	"github.com/Nostafa/golang-jwt/internal/store"
)

func main() {
	env.LoadEnv()

	cfg := config{
		port:    env.GetEnvString("PORT", ":8000"),
		env:     env.GetEnvString("ENV", "development"),
		version: env.GetEnvString("VERSION", "1.0.0"),
		db: dbConfig{
			url:          env.GetEnvString("DATABASE_URL", "postgres://postgres:v4yn96xklq9j5o4sk3dg@142.132.190.222:5434/test?sslmode=disable"),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
			maxLifetime:  env.GetEnvString("DB_MAX_LIFETIME", "10m"),
			maxIdleTime:  env.GetEnvString("DB_MAX_IDLE_TIME", "5m"),
		},
	}

	db, err := database.
		New(cfg.db.url, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxLifetime, cfg.db.maxIdleTime)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	storeDB := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  storeDB,
	}

	mux := app.mount()
	if err := app.server(mux); err != nil {
		log.Fatal(err)
	}
}
