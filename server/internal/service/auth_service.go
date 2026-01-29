package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/model"
	"github.com/shamal-iroshan/notora/internal/pkg/crypto"
	"github.com/shamal-iroshan/notora/internal/pkg/jwt"
	"github.com/shamal-iroshan/notora/internal/repository"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	TokenRepo *repository.TokenRepository
	ResetRepo *repository.ResetRepository
	AppConfig *config.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
	resetRepo *repository.ResetRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
		ResetRepo: resetRepo,
		AppConfig: cfg,
	}
}

// -----------------------------------------------------------------------------
// REGISTER
// -----------------------------------------------------------------------------

func (s *AuthService) Register(email, password, name string) error {

	// Check unique email
	u, _ := s.UserRepo.FindByEmail(email)
	if u != nil {
		return fmt.Errorf("email already registered")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user salt (used for encrypted notes master key)
	saltBytes := make([]byte, 16)
	rand.Read(saltBytes)
	userSalt := hex.EncodeToString(saltBytes)

	// Create new user (PENDING status by default)
	_, err = s.UserRepo.Create(
		email,
		string(passwordHash),
		name,
		userSalt,
		time.Now().UTC().Format(time.RFC3339),
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------
// LOGIN
// -----------------------------------------------------------------------------

func (s *AuthService) Login(email, password string) (string, string, error) {

	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Status checks
	if user.Status == "PENDING" {
		return "", "", fmt.Errorf("account not approved")
	}
	if user.Status == "SUSPENDED" {
		return "", "", fmt.Errorf("account suspended")
	}

	// Create access token
	accessToken, err := jwt.CreateAccess([]byte(s.AppConfig.JWTSecret), user.ID, s.AppConfig.AccessExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token
	refreshToken, err := crypto.RandomHex(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token")
	}

	err = s.TokenRepo.Insert(
		user.ID,
		crypto.SHA256Hex(refreshToken),
		time.Now().Add(time.Duration(s.AppConfig.RefreshExpiry)*time.Second),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// -----------------------------------------------------------------------------
// REFRESH TOKEN
// -----------------------------------------------------------------------------

func (s *AuthService) Refresh(oldRefreshToken string) (string, string, error) {

	hash := crypto.SHA256Hex(oldRefreshToken)

	tokenID, userID, err := s.TokenRepo.FindValid(hash)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// Revoke old token (rotation)
	if err := s.TokenRepo.Revoke(tokenID); err != nil {
		return "", "", fmt.Errorf("failed to revoke token: %w", err)
	}

	// Create new refresh token
	newRefresh, err := crypto.RandomHex(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token")
	}

	err = s.TokenRepo.Insert(
		userID,
		crypto.SHA256Hex(newRefresh),
		time.Now().Add(time.Duration(s.AppConfig.RefreshExpiry)*time.Second),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Create new access token
	newAccess, err := jwt.CreateAccess([]byte(s.AppConfig.JWTSecret), userID, s.AppConfig.AccessExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to create new access token")
	}

	return newAccess, newRefresh, nil
}

// -----------------------------------------------------------------------------
// FORGOT PASSWORD
// -----------------------------------------------------------------------------

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil // Always return OK
	}

	resetToken, _ := crypto.RandomHex(32)
	hash := crypto.SHA256Hex(resetToken)
	exp := time.Now().Add(10 * time.Minute)

	s.ResetRepo.Insert(user.ID, hash, exp)

	// In production send by SMTP
	fmt.Println("RESET URL: /reset-password?token=" + resetToken)

	return nil
}

// -----------------------------------------------------------------------------
// RESET PASSWORD
// -----------------------------------------------------------------------------

func (s *AuthService) ResetPassword(rawToken, newPassword string) error {

	hash := crypto.SHA256Hex(rawToken)

	resetID, userID, err := s.ResetRepo.FindValid(hash)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err := s.UserRepo.UpdatePassword(userID, string(newHash)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens (security)
	_ = s.TokenRepo.RevokeAllForUser(userID)

	s.ResetRepo.MarkUsed(resetID)

	return nil
}

// -----------------------------------------------------------------------------
// PROFILE EDIT
// -----------------------------------------------------------------------------

func (s *AuthService) EditProfile(userID int64, name string) error {
	return s.UserRepo.UpdateName(userID, name)
}

// -----------------------------------------------------------------------------
// CHANGE PASSWORD
// -----------------------------------------------------------------------------

func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return fmt.Errorf("old password incorrect")
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err := s.UserRepo.UpdatePassword(userID, string(newHash)); err != nil {
		return fmt.Errorf("failed to update password")
	}

	// Revoke all refresh tokens
	s.TokenRepo.RevokeAllForUser(userID)

	return nil
}

// -----------------------------------------------------------------------------
// GET USER
// -----------------------------------------------------------------------------

func (s *AuthService) GetUserByID(id int64) (*model.User, error) {
	return s.UserRepo.FindByID(id)
}
