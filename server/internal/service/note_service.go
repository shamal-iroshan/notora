package service

import (
	"database/sql"
	"errors"

	"github.com/shamal-iroshan/notora/internal/model"
	"github.com/shamal-iroshan/notora/internal/repository"
)

type NoteService struct {
	Repo *repository.NoteRepository
}

func NewNoteService(repo *repository.NoteRepository) *NoteService {
	return &NoteService{Repo: repo}
}

func (s *NoteService) Create(userID int64, title, content string) (int64, error) {
	return s.Repo.Create(userID, title, content)
}

func (s *NoteService) Get(userID, noteID int64) (interface{}, error) {
	note, err := s.Repo.GetByID(userID, noteID)
	if err != nil {
		return nil, errors.New("note not found")
	}
	return note, nil
}

func (s *NoteService) GetAll(userID int64) (*sql.Rows, error) {
	return s.Repo.GetAll(userID)
}

func (s *NoteService) Update(userID, noteID int64, title, content string) error {
	return s.Repo.Update(noteID, userID, title, content)
}

func (s *NoteService) UpdateFlags(userID, noteID int64, pinned, archived, deleted *bool) error {
	return s.Repo.UpdateFlags(noteID, userID, pinned, archived, deleted)
}

func (s *NoteService) DeleteForever(userID, noteID int64) error {
	return s.Repo.DeletePermanently(noteID, userID)
}

func (s *NoteService) Duplicate(userID, noteID int64) (int64, error) {
	return s.Repo.Duplicate(userID, noteID)
}

func (s *NoteService) Metadata(userID int64) (*sql.Rows, error) {
	return s.Repo.GetMetadata(userID)
}

func (s *NoteService) Search(userID int64, query string) ([]model.Note, error) {
	return s.Repo.Search(userID, query)
}
