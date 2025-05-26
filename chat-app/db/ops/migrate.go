package dbops

import (
	"chat-app/server/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)

func Migrate() error {
	// estabilish Connection
	c := config.AppConfig
	db, err := sql.Open("postgres", c.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// test connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}

	// read schema file
	schema, err := os.ReadFile("db/ops/schema.sql")
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
