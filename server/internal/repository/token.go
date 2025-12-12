package repository

import (
	"database/sql"
	"errors"
	"time"
)

// TokenRepository provides DB operations for managing refresh tokens.
// This includes inserting new tokens, verifying token validity,
// and revoking existing tokens.
type TokenRepository struct {
	DB *sql.DB
}

// NewTokenRepository creates a new instance of TokenRepository.
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

// Insert stores a new hashed refresh token for a specific user.
// The raw token is never stored—only its SHA-256 hash.
// This enhances security by preventing token leakage from the DB.
//
// uid: User ID that owns the token
// hash: SHA-256 hash of the refresh token
// exp: Expiration time for the token
func (r *TokenRepository) Insert(userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.DB.Exec(
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES (?, ?, ?)`,
		userID,
		tokenHash,
		expiresAt.UTC().Format(time.RFC3339),
	)

	return err
}

// FindValid retrieves a non-revoked, non-expired refresh token by its hash.
// It returns:
//   - tokenID: the token record ID
//   - userID: the owner of the token
//   - error: sql.ErrNoRows if token doesn't exist, is revoked, or expired
//
// This ensures invalid or expired tokens cannot be reused.
func (r *TokenRepository) FindValid(tokenHash string) (tokenID int64, userID int64, err error) {
	var revoked int64
	var expiresAtStr string

	err = r.DB.QueryRow(
		`SELECT id, user_id, expires_at, revoked
		   FROM refresh_tokens
		  WHERE token_hash = ?`,
		tokenHash,
	).Scan(&tokenID, &userID, &expiresAtStr, &revoked)

	if err != nil {
		return 0, 0, err // Could be sql.ErrNoRows → caller decides
	}

	// Reject revoked tokens
	if revoked != 0 {
		return 0, 0, sql.ErrNoRows
	}

	// Parse expiration timestamp
	expiresAt, parseErr := time.Parse(time.RFC3339, expiresAtStr)
	if parseErr != nil {
		return 0, 0, parseErr
	}

	// Reject expired tokens
	if time.Now().After(expiresAt) {
		return 0, 0, sql.ErrNoRows
	}

	return tokenID, userID, nil
}

// Revoke marks a refresh token as invalid so it cannot be used again.
// This is called when rotating tokens or during logout.
func (r *TokenRepository) Revoke(tokenID int64) error {
	result, err := r.DB.Exec(
		`UPDATE refresh_tokens
		    SET revoked = 1
		  WHERE id = ?`,
		tokenID,
	)

	if err != nil {
		return err
	}

	// Ensure at least 1 row was affected
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("no token found to revoke")
	}

	return nil
}
