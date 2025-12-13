package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/model"
	"github.com/shamal-iroshan/notora/internal/pkg/encryption"
)

type NoteRepository struct {
	DB        *sql.DB
	AppConfig *config.Config
}

func NewNoteRepository(db *sql.DB, cfg *config.Config) *NoteRepository {
	return &NoteRepository{DB: db, AppConfig: cfg}
}

func (r *NoteRepository) Create(userID int64, title, content string) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	encContent, err := encryption.EncryptAES([]byte(r.AppConfig.EncryptionKey), content)
	if err != nil {
		return 0, err
	}

	result, err := r.DB.Exec(`
		INSERT INTO notes (user_id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID, title, encContent, now, now)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *NoteRepository) GetByID(userID, noteID int64) (*model.Note, error) {
	var (
		id        int64
		title     string
		content   string
		pinned    int
		archived  int
		deleted   int
		createdAt string
		updatedAt string
	)

	err := r.DB.QueryRow(`
		SELECT id, title, content, is_pinned, is_archived, is_deleted, created_at, updated_at
		FROM notes
		WHERE user_id = ? AND id = ?
	`, userID, noteID).Scan(
		&id, &title, &content, &pinned, &archived, &deleted, &createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	plaintext, err := encryption.DecryptAES([]byte(r.AppConfig.EncryptionKey), content)
	if err != nil {
		return nil, err
	}

	return &model.Note{
		ID:         id,
		Title:      title,
		Content:    plaintext,
		IsPinned:   pinned == 1,
		IsArchived: archived == 1,
		IsDeleted:  deleted == 1,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func (r *NoteRepository) GetAll(userID int64) (*sql.Rows, error) {
	return r.DB.Query(`
		SELECT id, title, content, is_pinned, is_archived, is_deleted, created_at, updated_at
		FROM notes
		WHERE user_id = ?
		ORDER BY updated_at DESC
	`, userID)
}

func (r *NoteRepository) Update(noteID, userID int64, title, content string) error {
	// Encrypt content
	encContent, err := encryption.EncryptAES([]byte(r.AppConfig.EncryptionKey), content)
	if err != nil {
		return err
	}

	// Execute update query
	_, err = r.DB.Exec(`
        UPDATE notes
        SET title = ?, content = ?, updated_at = ?
        WHERE id = ? AND user_id = ?
    `, title, encContent, time.Now().UTC().Format(time.RFC3339), noteID, userID)

	return err
}

func (r *NoteRepository) UpdateFlags(noteID, userID int64, pinned, archived, deleted *bool) error {
	query := `UPDATE notes SET `
	args := []interface{}{}

	if pinned != nil {
		query += "is_pinned = ?, "
		args = append(args, boolToInt(*pinned))
	}

	if archived != nil {
		query += "is_archived = ?, "
		args = append(args, boolToInt(*archived))
	}

	if deleted != nil {
		query += "is_deleted = ?, "
		args = append(args, boolToInt(*deleted))
	}

	query += "updated_at = ? WHERE id = ? AND user_id = ?"
	args = append(args, time.Now().UTC().Format(time.RFC3339), noteID, userID)

	_, err := r.DB.Exec(query, args...)
	return err
}

func (r *NoteRepository) DeletePermanently(noteID, userID int64) error {
	_, err := r.DB.Exec(`DELETE FROM notes WHERE id = ? AND user_id = ?`, noteID, userID)
	return err
}

func (r *NoteRepository) Duplicate(userID, noteID int64) (int64, error) {
	var title, content string
	err := r.DB.QueryRow(`
		SELECT title, content
		FROM notes
		WHERE id = ? AND user_id = ?
	`, noteID, userID).Scan(&title, &content)

	if err != nil {
		return 0, err
	}

	now := time.Now().UTC().Format(time.RFC3339)

	result, err := r.DB.Exec(`
		INSERT INTO notes (user_id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID, title+" (Copy)", content, now, now)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *NoteRepository) GetMetadata(userID int64) (*sql.Rows, error) {
	return r.DB.Query(`
		SELECT id, title, updated_at, is_pinned, is_archived, is_deleted
		FROM notes
		WHERE user_id = ?
		ORDER BY updated_at DESC
	`, userID)
}

func (r *NoteRepository) Search(userID int64, query string) ([]model.Note, error) {
	query = "%" + query + "%" // wildcard search

	rows, err := r.DB.Query(`
		SELECT id, title, content, is_pinned, is_archived, is_deleted, created_at, updated_at
		FROM notes
		WHERE user_id = ?
		  AND is_deleted = 0
		  AND (title LIKE ? OR content LIKE ?)
		ORDER BY updated_at DESC
	`, userID, query, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Note

	for rows.Next() {
		var n model.Note
		var pinned, archived, deleted int

		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&pinned,
			&archived,
			&deleted,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}

		n.IsPinned = pinned == 1
		n.IsArchived = archived == 1
		n.IsDeleted = deleted == 1

		results = append(results, n)
	}

	return results, nil
}

// EnsureOwnership verifies that the note belongs to the given user.
// Returns nil if the user owns the note, otherwise error.
func (r *NoteRepository) EnsureOwnership(userID, noteID int64) error {
	var count int

	err := r.DB.QueryRow(`
		SELECT COUNT(1)
		FROM notes
		WHERE id = ? AND user_id = ?
	`, noteID, userID).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("not owner")
	}

	return nil
}

func (r *NoteRepository) GetPublicNote(noteID int64) (*model.Note, error) {
	var n model.Note
	var encContent string
	var pinned, archived, deleted int

	err := r.DB.QueryRow(`
		SELECT id, title, content, is_pinned, is_archived, is_deleted, created_at, updated_at
		FROM notes
		WHERE id = ?
	`, noteID).Scan(
		&n.ID,
		&n.Title,
		&encContent,
		&pinned,
		&archived,
		&deleted,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Decrypt content
	plaintext, err := encryption.DecryptAES([]byte(r.AppConfig.EncryptionKey), encContent)
	if err != nil {
		return nil, err
	}

	n.Content = plaintext
	n.IsPinned = pinned == 1
	n.IsArchived = archived == 1
	n.IsDeleted = deleted == 1

	return &n, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
