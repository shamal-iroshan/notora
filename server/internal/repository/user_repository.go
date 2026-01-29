package repository

import (
	"database/sql"

	"github.com/shamal-iroshan/notora/internal/model"
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
	res, err := r.DB.Exec(`INSERT INTO users(email, password_hash, name, user_salt, status, is_admin, created_at)
                           VALUES (?, ?, ?, ?, 'PENDING', 0, ?)`,
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
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var u model.User

	err := r.DB.QueryRow(`SELECT id, password_hash, email, name, user_salt, status, is_admin,created_at 
                         FROM users WHERE email = ?`,
		email).Scan(&u.ID, &u.Password, &u.Email, &u.Name, &u.UserSalt, &u.Status, &u.IsAdmin, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID retrieves user fields by id.
//
// Returns: id, email, passwordHash, name, err
func (r *UserRepository) FindByID(userID int64) (*model.User, error) {
	var u model.User

	err := r.DB.QueryRow(`
		SELECT id, email, password_hash, name, user_salt, status, is_admin, created_at
		FROM users WHERE id=?
	`, userID).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.UserSalt, &u.Status, &u.IsAdmin, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateName(userID int64, name string) error {
	_, err := r.DB.Exec(`UPDATE users SET name = ? WHERE id = ?`, name, userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID int64, newHash string) error {
	_, err := r.DB.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, newHash, userID)
	return err
}

// ADMIN ACTIONS
func (r *UserRepository) Approve(id int64) error {
	_, err := r.DB.Exec(`UPDATE users SET status='APPROVED' WHERE id=?`, id)
	return err
}

func (r *UserRepository) Suspend(id int64) error {
	_, err := r.DB.Exec(`UPDATE users SET status='SUSPENDED' WHERE id=?`, id)
	return err
}

func (r *UserRepository) DeleteUser(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id=?`, id)
	return err
}

func (r *UserRepository) ListPending() ([]model.User, error) {
	rows, err := r.DB.Query(`SELECT id, email, name, created_at FROM users WHERE status='PENDING'`)
	if err != nil {
		return nil, err
	}

	var list []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
		list = append(list, u)
	}
	return list, nil
}
