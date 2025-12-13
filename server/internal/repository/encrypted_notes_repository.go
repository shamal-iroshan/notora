package repository

import (
	"database/sql"
	"time"

	"github.com/shamal-iroshan/notora/internal/model"
)

type EncryptedNotesRepository struct {
	DB *sql.DB
}

func NewEncryptedNotesRepository(db *sql.DB) *EncryptedNotesRepository {
	return &EncryptedNotesRepository{DB: db}
}

// Create a new encrypted note
func (r *EncryptedNotesRepository) Create(userID int64, title, content, tnonce, cnonce, salt string) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	res, err := r.DB.Exec(`
        INSERT INTO encrypted_notes 
        (user_id, title_ciphertext, content_ciphertext, title_nonce, content_nonce, note_salt, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, userID, title, content, tnonce, cnonce, salt, now, now)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// List encrypted notes metadata (includes encrypted title)
func (r *EncryptedNotesRepository) List(userID int64) ([]model.EncryptedNoteMetadata, error) {
	rows, err := r.DB.Query(`
        SELECT id, title_ciphertext, title_nonce, note_salt, created_at, updated_at
        FROM encrypted_notes
        WHERE user_id = ?
        ORDER BY updated_at DESC
    `, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []model.EncryptedNoteMetadata{}
	for rows.Next() {
		var n model.EncryptedNoteMetadata
		if err := rows.Scan(&n.ID, &n.TitleCiphertext, &n.TitleNonce, &n.NoteSalt, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, nil
}

// Get full encrypted note
func (r *EncryptedNotesRepository) GetByID(userID, noteID int64) (*model.EncryptedNoteResponse, error) {
	var n model.EncryptedNoteResponse

	err := r.DB.QueryRow(`
        SELECT id, title_ciphertext, content_ciphertext, title_nonce, content_nonce, note_salt, created_at, updated_at
        FROM encrypted_notes
        WHERE id = ? AND user_id = ?
    `, noteID, userID).Scan(
		&n.ID, &n.TitleCiphertext, &n.ContentCiphertext,
		&n.TitleNonce, &n.ContentNonce, &n.NoteSalt,
		&n.CreatedAt, &n.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

// Update encrypted note
func (r *EncryptedNotesRepository) Update(userID, noteID int64, title, content, tnonce, cnonce, salt string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := r.DB.Exec(`
        UPDATE encrypted_notes
        SET title_ciphertext = ?, content_ciphertext = ?, title_nonce = ?, content_nonce = ?, 
            note_salt = ?, updated_at = ?
        WHERE id = ? AND user_id = ?
    `, title, content, tnonce, cnonce, salt, now, noteID, userID)

	return err
}

// Delete encrypted note
func (r *EncryptedNotesRepository) Delete(userID, noteID int64) error {
	_, err := r.DB.Exec(`DELETE FROM encrypted_notes WHERE id = ? AND user_id = ?`, noteID, userID)
	return err
}

// Check if a note is encrypted (for blocking sharing)
func (r *EncryptedNotesRepository) IsEncryptedNote(noteID int64) (bool, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM encrypted_notes WHERE id = ?`, noteID).Scan(&count)
	return count > 0, err
}
