package main

import (
	dbops "chat-app/db/ops"
	"chat-app/server/config"
	h "chat-app/server/handlers/http"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize configuration
	config.Init()

	// Init db connection
	dbConn, err := sql.Open("postgres", config.AppConfig.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer dbConn.Close()

	// Initialize user queries with the db connection
	db := dbops.NewUserQueries(dbConn)

	// Run db migration
	if err := dbops.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize handlers
	authHandler := h.NewAuthHandler(db)
	userHandler := h.NewUserHandler(db)
	
	// Initialize router
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// User routes
	router.HandleFunc("/api/users/profile", userHandler.GetProfile).Methods("GET")
	router.HandleFunc("/api/users/profile", userHandler.UpdateProfile).Methods("PUT")
	router.HandleFunc("/api/users", userHandler.GetUsers).Methods("GET")
	
	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}