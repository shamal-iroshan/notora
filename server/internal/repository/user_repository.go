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
func (r *UserRepository) Create(email, hash, name, salt, created string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO users(email, password_hash, name, user_salt, created_at)
                           VALUES (?, ?, ?, ?, ?)`,
		email, hash, name, salt, created)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// FindByEmail retrieves a user record by email.
//
// Returns:
//   - id: user ID
//   - passwordHash: stored bcrypt hash for authentication
//   - name: the user’s name
//   - createdAt: account creation timestamp
//   - error: sql.ErrNoRows if user does not exist
func (r *UserRepository) FindByEmail(email string) (id int64, hash, name, salt string, created string, err error) {

	err = r.DB.QueryRow(`SELECT id, password_hash, name, user_salt, created_at 
                         FROM users WHERE email = ?`,
		email).Scan(&id, &hash, &name, &salt, &created)

	// err may be sql.ErrNoRows or something else — caller handles it.
	return
}

// FindByID retrieves user fields by id.
//
// Returns: id, email, passwordHash, name, err
func (r *UserRepository) FindByID(userID int64) (id int64, email, passwordHash, salt string, name string, created string, err error) {
	err = r.DB.QueryRow(
		`SELECT id, email, password_hash,, user_salt, name, created_at 
		   FROM users
		  WHERE id = ?`,
		userID,
	).Scan(&id, &email, &passwordHash, &salt, &name, &created)
	return
}

func (r *UserRepository) UpdateName(userID int64, name string) error {
	_, err := r.DB.Exec(`UPDATE users SET name = ? WHERE id = ?`, name, userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID int64, newHash string) error {
	_, err := r.DB.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, newHash, userID)
	return err
}
