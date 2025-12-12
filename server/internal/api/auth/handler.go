package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/repository"
	"github.com/shamal-iroshan/notora/internal/service"
)

// AuthHandler holds dependencies for authentication routes.
// It connects the HTTP layer (Gin) → Service layer → Repository layer.
type AuthHandler struct {
	AuthService *service.AuthService
	AppConfig   *config.Config
}

// NewAuthHandler wires repositories → services → handler.
// This follows a clean architecture dependency order.
func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	authService := service.NewAuthService(userRepo, tokenRepo, cfg)

	return &AuthHandler{
		AuthService: authService,
		AppConfig:   cfg,
	}
}

// setCookie wraps Gin's cookie setter and ensures consistent configuration.
// HttpOnly prevents JavaScript access. Secure depends on production mode.
func setCookie(ctx *gin.Context, name, value string, maxAgeSeconds int, cfg *config.Config) {
	ctx.SetCookie(
		name,
		value,
		maxAgeSeconds,
		"/",
		cfg.CookieDomain,
		cfg.CookieSecure,
		true, // HttpOnly
	)
}

// -----------------------------------------------------------------------------
// REGISTER
// -----------------------------------------------------------------------------

// Register handles creation of a new user account.
func (h *AuthHandler) Register(ctx *gin.Context) {
	// Request body structure
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	// Validate incoming request
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service layer to perform validation + DB insert
	if err := h.AuthService.Register(requestBody.Email, requestBody.Password, requestBody.Name); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "account_created"})
}

// -----------------------------------------------------------------------------
// LOGIN
// -----------------------------------------------------------------------------

// Login validates user credentials and issues access + refresh cookies.
func (h *AuthHandler) Login(ctx *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Access token and refresh token returned from service
	accessToken, refreshToken, err := h.AuthService.Login(requestBody.Email, requestBody.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Set secure HttpOnly cookies
	setCookie(ctx, "access_token", accessToken, h.AppConfig.AccessExpiry, h.AppConfig)
	setCookie(ctx, "refresh_token", refreshToken, h.AppConfig.RefreshExpiry, h.AppConfig)

	ctx.JSON(http.StatusOK, gin.H{"status": "logged_in"})
}

// -----------------------------------------------------------------------------
// REFRESH TOKEN
// -----------------------------------------------------------------------------

// Refresh issues a new access token + refresh token using the existing refresh token.
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	newAccessToken, newRefreshToken, err := h.AuthService.Refresh(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	setCookie(ctx, "access_token", newAccessToken, h.AppConfig.AccessExpiry, h.AppConfig)
	setCookie(ctx, "refresh_token", newRefreshToken, h.AppConfig.RefreshExpiry, h.AppConfig)

	ctx.JSON(http.StatusOK, gin.H{"status": "token_refreshed"})
}

// -----------------------------------------------------------------------------
// GET CURRENT USER
// -----------------------------------------------------------------------------

// Me returns basic information about the currently logged-in user.
// The user_id is injected into the context by JWT middleware.
func (h *AuthHandler) Me(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

// -----------------------------------------------------------------------------
// LOGOUT
// -----------------------------------------------------------------------------

// Logout removes all auth cookies. Refresh tokens in the DB can also be revoked.
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// Immediately expire both cookies
	setCookie(ctx, "access_token", "", -1, h.AppConfig)
	setCookie(ctx, "refresh_token", "", -1, h.AppConfig)

	ctx.JSON(http.StatusOK, gin.H{"status": "logged_out"})
}
