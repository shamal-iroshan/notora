package service

import (
	"fmt"

	"github.com/shamal-iroshan/notora/internal/model"
	"github.com/shamal-iroshan/notora/internal/pkg/crypto"
	"github.com/shamal-iroshan/notora/internal/repository"
)

type ShareService struct {
	Notes *repository.NoteRepository
	Share *repository.ShareRepository
}

func NewShareService(n *repository.NoteRepository, s *repository.ShareRepository) *ShareService {
	return &ShareService{Notes: n, Share: s}
}

func (s *ShareService) CreateShare(userID, noteID int64) (string, error) {
	// Make sure user owns note
	if err := s.Notes.EnsureOwnership(userID, noteID); err != nil {
		return "", fmt.Errorf("not allowed")
	}

	token, _ := crypto.RandomHex(32)

	if err := s.Share.Create(noteID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *ShareService) DisableShare(userID, noteID int64) error {
	if err := s.Notes.EnsureOwnership(userID, noteID); err != nil {
		return fmt.Errorf("not allowed")
	}

	return s.Share.Disable(noteID)
}

func (s *ShareService) GetSharedNote(token string) (*model.Note, error) {
	noteID, err := s.Share.FindByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// Get note (service decrypts via repository)
	return s.Notes.GetPublicNote(noteID)
}
