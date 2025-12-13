package auth

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes registers routes that do NOT require authentication.
// These cover user registration, login, refresh token rotation,
// and the forgot-password / reset-password flow.
func RegisterPublicRoutes(router *gin.RouterGroup, handler *AuthHandler) {

	// POST /api/auth/register → Create a new user account
	router.POST("/register", handler.Register)

	// POST /api/auth/login → Login and set cookies for access + refresh tokens
	router.POST("/login", handler.Login)

	// POST /api/auth/refresh → Issue new access & refresh tokens
	router.POST("/refresh", handler.Refresh)

	// POST /api/auth/forgot-password → Start password reset process
	router.POST("/forgot-password", handler.ForgotPassword)

	// POST /api/auth/reset-password → Complete password reset using token
	router.POST("/reset-password", handler.ResetPassword)
}

// RegisterProtectedRoutes registers routes that REQUIRE authentication.
// These routes can only be accessed with a valid JWT access token cookie.
func RegisterProtectedRoutes(router *gin.RouterGroup, handler *AuthHandler, authMiddleware gin.HandlerFunc) {

	// Apply JWT middleware to all protected endpoints
	router.Use(authMiddleware)

	// GET /api/me → Get authenticated user's info
	router.GET("/me", handler.Me)

	// PUT /api/me → Update user profile (e.g., name)
	router.PUT("/me", handler.EditProfile)

	// PUT /api/me/password → Change password (requires old password)
	router.PUT("/me/password", handler.ChangePassword)

	// POST /api/logout → Invalidate refresh tokens & clear cookies
	router.POST("/logout", handler.Logout)
}
