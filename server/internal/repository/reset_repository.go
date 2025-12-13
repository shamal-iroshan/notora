package repository

import (
	"database/sql"
	"time"
)

type ResetRepository struct {
	DB *sql.DB
}

func NewResetRepository(db *sql.DB) *ResetRepository {
	return &ResetRepository{DB: db}
}

func (r *ResetRepository) Insert(userID int64, tokenHash string, expires time.Time) error {
	_, err := r.DB.Exec(`
		INSERT INTO password_resets (user_id, token_hash, expires_at, created_at)
		VALUES (?, ?, ?, ?)
	`, userID, tokenHash, expires.UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	return err
}

func (r *ResetRepository) FindValid(tokenHash string) (id, userID int64, err error) {
	var expiresStr string
	var used int64

	err = r.DB.QueryRow(`
		SELECT id, user_id, expires_at, used
		FROM password_resets
		WHERE token_hash = ?
	`, tokenHash).Scan(&id, &userID, &expiresStr, &used)

	if err != nil {
		return
	}

	if used != 0 {
		return 0, 0, sql.ErrNoRows
	}

	exp, _ := time.Parse(time.RFC3339, expiresStr)
	if time.Now().After(exp) {
		return 0, 0, sql.ErrNoRows
	}

	return
}

func (r *ResetRepository) MarkUsed(id int64) error {
	_, err := r.DB.Exec(`UPDATE password_resets SET used = 1 WHERE id = ?`, id)
	return err
}
