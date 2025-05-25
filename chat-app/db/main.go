package main

import (
	"log"
	"chat-app/db/ops"
)

func main() {
	cfg := &db.Config{
		DB_HOST: "localhost",
		DB_PORT: 5432,
		DB_USER: "postgres",
		DB_PASSWORD: "postgres",
		DB_NAME: "chat-app",
		DB_SSLMODE: "disable",
	}

	if err := db.Migrate(cfg); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}