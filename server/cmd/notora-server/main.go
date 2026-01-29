package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	// Auth modules
	"github.com/shamal-iroshan/notora/internal/api/admin"
	"github.com/shamal-iroshan/notora/internal/api/auth"
	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/db"
	"github.com/shamal-iroshan/notora/internal/middleware"

	// Notes modules
	encryptedapi "github.com/shamal-iroshan/notora/internal/api/encrypted"
	noteapi "github.com/shamal-iroshan/notora/internal/api/notes"
	shareapi "github.com/shamal-iroshan/notora/internal/api/share"
	"github.com/shamal-iroshan/notora/internal/repository"
	"github.com/shamal-iroshan/notora/internal/service"
)

func main() {
	// -------------------------------------------------------------
	// Load environment variables (.env file in development)
	// -------------------------------------------------------------
	_ = godotenv.Load()

	// Load application configuration (port, DB path, secrets, cookies)
	cfg := config.LoadFromEnv()

	// -------------------------------------------------------------
	// Ensure data directory exists (for SQLite DB)
	// -------------------------------------------------------------
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatal("failed to create data directory:", err)
	}

	// -------------------------------------------------------------
	// Connect to SQLite database (enable foreign key constraints)
	// -------------------------------------------------------------
	dbConn, err := sql.Open("sqlite3", cfg.DBPath+"?_fk=1")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer dbConn.Close()

	// -------------------------------------------------------------
	// Database migrations (create tables if not exist)
	// -------------------------------------------------------------
	if err := db.Migrate(dbConn); err != nil {
		log.Fatal("database migration failed:", err)
	}

	// -------------------------------------------------------------
	// Initialize Gin HTTP server
	// -------------------------------------------------------------
	r := gin.Default()

	// middlewares
	userRepo := repository.NewUserRepository(dbConn)

	pendingBlock := middleware.RequireApprovedUser()
	jwtBlock := middleware.JWTMiddleware(cfg, userRepo)

	// -------------------------------------------------------------
	// AUTHENTICATION SETUP
	// -------------------------------------------------------------
	authHandler := auth.NewAuthHandler(dbConn, cfg)

	// Public auth routes (register, login, refresh, forgot/reset password)
	auth.RegisterPublicRoutes(r.Group("/api/auth"), authHandler)

	// Protected auth routes (me, edit profile, change password, logout)
	auth.RegisterProtectedRoutes(
		r.Group("/api"),
		authHandler,
		jwtBlock,
	)

	// -------------------------------------------------------------
	// NOTES MODULE SETUP
	// -------------------------------------------------------------

	// Create Note repository → service → handler
	noteRepo := repository.NewNoteRepository(dbConn, cfg)
	noteService := service.NewNoteService(noteRepo)
	noteHandler := noteapi.NewNoteHandler(noteService)

	// Register protected notes routes
	// All note endpoints require authentication.
	noteapi.RegisterNoteRoutes(
		r.Group("/api"),
		noteHandler,
		jwtBlock,
		pendingBlock,
	)

	// -------------------------------
	// SHARING MODULE SETUP
	// -------------------------------
	shareRepo := repository.NewShareRepository(dbConn)
	shareService := service.NewShareService(noteRepo, shareRepo)
	shareHandler := shareapi.NewShareHandler(shareService, cfg)

	// Public sharing: no auth
	shareapi.RegisterPublicShareRoutes(r.Group("/api"), shareHandler)

	// Protected sharing: must own the note
	shareapi.RegisterProtectedShareRoutes(r.Group("/api", jwtBlock, pendingBlock), shareHandler)

	encryptedRepo := repository.NewEncryptedNotesRepository(dbConn)
	encryptedService := service.NewEncryptedNotesService(encryptedRepo)
	encryptedHandler := encryptedapi.NewEncryptedNotesHandler(encryptedService)

	encryptedapi.RegisterEncryptedNotesRoutes(
		r.Group("/api/encrypted-notes", jwtBlock, pendingBlock),
		encryptedHandler,
	)

	// Admin Area
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(jwtBlock, middleware.RequireAdmin())
	admin.RegisterAdminRoutes(adminGroup, adminHandler)

	// -------------------------------------------------------------
	// START SERVER
	// Example: :8000
	// -------------------------------------------------------------
	log.Println("NOTORA server running on port:", cfg.Port)
	r.Run(":" + cfg.Port)
}
