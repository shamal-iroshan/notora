package repository

import (
	"database/sql"
	"time"
)

type ShareRepository struct {
	DB *sql.DB
}

func NewShareRepository(db *sql.DB) *ShareRepository {
	return &ShareRepository{DB: db}
}

func (r *ShareRepository) Create(noteID int64, token string) error {
	_, err := r.DB.Exec(`
		INSERT INTO shared_notes (note_id, token, created_at)
		VALUES (?, ?, ?)
	`, noteID, token, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (r *ShareRepository) FindByToken(token string) (int64, error) {
	var noteID int64
	var disabled int

	err := r.DB.QueryRow(`
		SELECT note_id, disabled
		FROM shared_notes
		WHERE token = ?
	`, token).Scan(&noteID, &disabled)

	if err != nil {
		return 0, err
	}
	if disabled != 0 {
		return 0, sql.ErrNoRows
	}

	return noteID, nil
}

func (r *ShareRepository) Disable(noteID int64) error {
	_, err := r.DB.Exec(`
		UPDATE shared_notes
		SET disabled = 1
		WHERE note_id = ?
	`, noteID)
	return err
}
