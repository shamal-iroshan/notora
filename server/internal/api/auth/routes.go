package auth

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes registers routes that do NOT require authentication.
// These endpoints handle creating accounts, logging in, and refreshing tokens.
func RegisterPublicRoutes(router *gin.RouterGroup, handler *AuthHandler) {
	// POST /api/auth/register → Create new account
	router.POST("/register", handler.Register)

	// POST /api/auth/login → Login and issue access + refresh tokens
	router.POST("/login", handler.Login)

	// POST /api/auth/refresh → Refresh access token using refresh cookie
	router.POST("/refresh", handler.Refresh)
}

// RegisterProtectedRoutes registers routes that REQUIRE authentication.
// The provided middleware (mw) is applied to the entire group.
func RegisterProtectedRoutes(router *gin.RouterGroup, handler *AuthHandler, authMiddleware gin.HandlerFunc) {
	// Apply JWT authentication middleware to all protected endpoints
	router.Use(authMiddleware)

	// GET /api/me → Returns authenticated user's info
	router.GET("/me", handler.Me)

	// POST /api/logout → Clears cookies and invalidates refresh token
	router.POST("/logout", handler.Logout)
}
