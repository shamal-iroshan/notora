package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/pkg/crypto"
	"github.com/shamal-iroshan/notora/internal/pkg/jwt"
	"github.com/shamal-iroshan/notora/internal/repository"
)

// AuthService contains the business logic for authentication.
// It coordinates user repository + token repository + JWT utilities.
type AuthService struct {
	UserRepo  *repository.UserRepository
	TokenRepo *repository.TokenRepository
	ResetRepo *repository.ResetRepository
	AppConfig *config.Config
}

// NewAuthService wires dependencies into the service layer.
func NewAuthService(
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
		AppConfig: cfg,
	}
}

// -----------------------------------------------------------------------------
// REGISTER
// -----------------------------------------------------------------------------

// Register creates a new user account after verifying that the email is unique.
func (s *AuthService) Register(email, password, name string) error {
	// Check if user already exists
	_, _, _, _, _, err := s.UserRepo.FindByEmail(email)
	if err == nil {
		return fmt.Errorf("email already registered")
	}

	// Hash password using bcrypt (safe for storing passwords)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	//  Generate user_salt (for master key derivation)
	saltBytes := make([]byte, 16)
	_, _ = rand.Read(saltBytes)
	userSalt := hex.EncodeToString(saltBytes)

	// Insert new user
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

// Login verifies email + password, then issues access & refresh tokens.
func (s *AuthService) Login(email, password string) (accessToken string, refreshToken string, err error) {

	// Retrieve user by email
	userID, storedHash, _, _, _, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Verify password
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Create JWT access token
	accessToken, err = jwt.CreateAccess([]byte(s.AppConfig.JWTSecret), userID, s.AppConfig.AccessExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token (raw token stored in cookie, hash stored in DB)
	refreshToken, err = crypto.RandomHex(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.TokenRepo.Insert(
		userID,
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

// Refresh handles rotation of refresh tokens and issues a new access token.
func (s *AuthService) Refresh(oldRefreshToken string) (newAccessToken string, newRefreshToken string, err error) {

	hashedToken := crypto.SHA256Hex(oldRefreshToken)

	// Validate existing refresh token
	refreshTokenID, userID, err := s.TokenRepo.FindValid(hashedToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// Revoke old refresh token (refresh token rotation)
	if err := s.TokenRepo.Revoke(refreshTokenID); err != nil {
		return "", "", fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err = crypto.RandomHex(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	err = s.TokenRepo.Insert(
		userID,
		crypto.SHA256Hex(newRefreshToken),
		time.Now().Add(time.Duration(s.AppConfig.RefreshExpiry)*time.Second),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to store new refresh token: %w", err)
	}

	// Issue new access token
	newAccessToken, err = jwt.CreateAccess([]byte(s.AppConfig.JWTSecret), userID, s.AppConfig.AccessExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to create new access token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	// check user exists
	userID, _, _, _, _, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil // always return OK for security
	}

	// create reset token
	resetToken, _ := crypto.RandomHex(32)
	tokenHash := crypto.SHA256Hex(resetToken)

	expiry := time.Now().Add(10 * time.Minute)

	s.ResetRepo.Insert(userID, tokenHash, expiry)

	// In production you send this via SMTP
	fmt.Println("RESET LINK: /reset-password?token=" + resetToken)

	return nil
}

func (s *AuthService) ResetPassword(rawToken, newPassword string) error {
	tokenHash := crypto.SHA256Hex(rawToken)

	resetID, userID, err := s.ResetRepo.FindValid(tokenHash)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	// inside AuthService.ResetPassword(...)
	if err := s.UserRepo.UpdatePassword(userID, string(newHash)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// revoke all refresh tokens for the user
	if err := s.TokenRepo.RevokeAllForUser(userID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}
	s.ResetRepo.MarkUsed(resetID)

	return nil
}

func (s *AuthService) EditProfile(userID int64, name string) error {
	return s.UserRepo.UpdateName(userID, name)
}

func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	// Fetch user data
	_, _, storedHash, _, _, _, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(oldPassword)) != nil {
		return fmt.Errorf("old password incorrect")
	}

	// Hash new password
	newHashBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password")
	}

	// Update database
	if err := s.UserRepo.UpdatePassword(userID, string(newHashBytes)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Security: revoke all refresh tokens
	_ = s.TokenRepo.RevokeAllForUser(userID)

	return nil
}

func (s *AuthService) GetUserByID(userID int64) (id int64, email, passwordHash, salt string, name string, created string, err error) {
	return s.UserRepo.FindByID(userID)
}
