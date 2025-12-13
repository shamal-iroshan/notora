package service

import (
	"github.com/shamal-iroshan/notora/internal/model"
	"github.com/shamal-iroshan/notora/internal/repository"
)

type EncryptedNotesService struct {
	Repo *repository.EncryptedNotesRepository
}

func NewEncryptedNotesService(r *repository.EncryptedNotesRepository) *EncryptedNotesService {
	return &EncryptedNotesService{Repo: r}
}

func (s *EncryptedNotesService) Create(userID int64, dto model.CreateEncryptedNoteInput) (int64, error) {
	return s.Repo.Create(
		userID,
		dto.TitleCiphertext,
		dto.ContentCiphertext,
		dto.TitleNonce,
		dto.ContentNonce,
		dto.NoteSalt,
	)
}

func (s *EncryptedNotesService) List(userID int64) ([]model.EncryptedNoteMetadata, error) {
	return s.Repo.List(userID)
}

func (s *EncryptedNotesService) Get(userID, noteID int64) (*model.EncryptedNoteResponse, error) {
	return s.Repo.GetByID(userID, noteID)
}

func (s *EncryptedNotesService) Update(userID, noteID int64, dto model.UpdateEncryptedNoteInput) error {
	return s.Repo.Update(
		userID,
		noteID,
		dto.TitleCiphertext,
		dto.ContentCiphertext,
		dto.TitleNonce,
		dto.ContentNonce,
		dto.NoteSalt,
	)
}

func (s *EncryptedNotesService) Delete(userID, noteID int64) error {
	return s.Repo.Delete(userID, noteID)
}

func (s *EncryptedNotesService) IsEncryptedNote(noteID int64) (bool, error) {
	return s.Repo.IsEncryptedNote(noteID)
}
