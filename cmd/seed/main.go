package main

import (
	"log"

	"github.com/Nostafa/golang-jwt/internal/database"
	"github.com/Nostafa/golang-jwt/internal/env"
	"github.com/Nostafa/golang-jwt/internal/store"
)

func main() {

	env.LoadEnv()

	db, err := database.New(
		env.GetEnvString("DATABASE_URL", "postgres://postgres:v4yn96xklq9j5o4sk3dg@142.132.190.222:5434/test?sslmode=disable"),
		env.GetEnvInt("DB_MAX_OPEN_CONNS", 10),
		env.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
		env.GetEnvString("DB_MAX_LIFETIME", "10m"),
		env.GetEnvString("DB_MAX_IDLE_TIME", "5m"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	storeDB := store.NewStorage(db)
	err = database.Seed(storeDB)
	if err != nil {
		log.Fatal(err)
	}
}
