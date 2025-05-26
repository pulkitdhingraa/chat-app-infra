package dbops

import (
	"database/sql"
	"errors"
)

type UserQueries struct {
	db *sql.DB
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
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

func (q *UserQueries) GetUserByID(id string) (*User, error) {
	user := &User{}
	err := q.db.QueryRow(`
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (q *UserQueries) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := q.db.QueryRow(
		"SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (q *UserQueries) UpdateUser(id, username, email string) error {
	_, err := q.db.Exec(`
		UPDATE users
		SET username = $1, email = $2
		WHERE id = $3`,
		username, email, id,
	)
	return err
}

func (q *UserQueries) GetUsers(currentUserID string) ([]User, error) {
	rows, err := q.db.Query(`
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id != $1
		ORDER BY username`,
		currentUserID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
