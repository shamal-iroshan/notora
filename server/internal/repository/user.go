package repository

import (
	"database/sql"
)

// UserRepository provides database operations for the users table.
// It handles user creation and lookup by email.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository returns an initialized UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create inserts a new user record into the database.
//
// Parameters:
//   - email: user's email (must be unique)
//   - passwordHash: bcrypt hashed password
//   - name: optional display name
//   - createdAt: timestamp when the user was created (UTC, RFC3339)
//
// Returns:
//   - userID: the auto-incremented ID of the new user
//   - error, if any database error occurred (e.g., email already exists)
func (r *UserRepository) Create(email, passwordHash, name, createdAt string) (int64, error) {
	result, err := r.DB.Exec(
		`INSERT INTO users (email, password_hash, name, created_at)
		 VALUES (?, ?, ?, ?)`,
		email, passwordHash, name, createdAt,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindByEmail retrieves a user record by email.
//
// Returns:
//   - id: user ID
//   - passwordHash: stored bcrypt hash for authentication
//   - name: the user’s name
//   - createdAt: account creation timestamp
//   - error: sql.ErrNoRows if user does not exist
func (r *UserRepository) FindByEmail(email string) (id int64, passwordHash, name, createdAt string, err error) {

	err = r.DB.QueryRow(
		`SELECT id, password_hash, name, created_at
		   FROM users
		  WHERE email = ?`,
		email,
	).Scan(&id, &passwordHash, &name, &createdAt)

	// err may be sql.ErrNoRows or something else — caller handles it.
	return
}
