package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSLMODE  string
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.DB_HOST, c.DB_PORT, c.DB_USER, c.DB_PASSWORD, c.DB_NAME, c.DB_SSLMODE)
}

func Migrate(cfg *Config) error {
	// estabilish Connection
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// test connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}

	// read schema file
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	// execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("error executing schema: %v", err)
	}

	log.Println("Schema migration completed successfully")
	return nil
}
