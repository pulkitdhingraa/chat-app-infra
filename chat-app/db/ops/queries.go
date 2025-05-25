package db

import (
	"database/sql"
)

type UserQueries struct {
	db *sql.DB
}

func NewUserQueries(db *sql.DB) *UserQueries {
	return &UserQueries{db: db}
}

func (q *UserQueries) CheckUserExists(email, username string) (bool, error) {
	var exists bool
	err := q.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)", email, username,
	).Scan(&exists)
	return exists, err
}

func (q *UserQueries) CreateUser(username, email, password_hash string) (string, error) {
	var userID string
	err := q.db.QueryRow(`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id`,
		username, email, password_hash,
	).Scan(&userID)
	return userID, err
}