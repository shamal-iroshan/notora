package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/shamal-iroshan/notora/internal/api/auth"
	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/db"
	"github.com/shamal-iroshan/notora/internal/middleware"
)

func main() {
	// -------------------------------------------------------------
	// Load environment variables from .env (only in development)
	// This does NOT override system environment variables.
	// -------------------------------------------------------------
	_ = godotenv.Load()

	// Load application configuration (port, DB path, secrets, etc.)
	cfg := config.LoadFromEnv()

	// -------------------------------------------------------------
	// Ensure the data directory exists (for SQLite .db file)
	// If the folder does not exist, it will be created.
	// -------------------------------------------------------------
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------
	// Connect to SQLite database
	// "_fk=1" enables foreign key constraints in SQLite.
	// -------------------------------------------------------------
	dbConn, err := sql.Open("sqlite3", cfg.DBPath+"?_fk=1")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// -------------------------------------------------------------
	// Run database migrations (create tables if they don't exist)
	// This ensures your database schema is ready before the server starts.
	// -------------------------------------------------------------
	if err := db.Migrate(dbConn); err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------
	// Initialize Gin web server
	// gin.Default() includes logger + recovery middleware
	// -------------------------------------------------------------
	r := gin.Default()

	// Create authentication handler with its dependencies (DB + config)
	authHandler := auth.NewAuthHandler(dbConn, cfg)

	// -------------------------------------------------------------
	// Register API routes
	//
	// /api/auth → Public routes (login, register, refresh)
	// /api      → Protected routes that require JWT (me, logout, etc.)
	// -------------------------------------------------------------
	auth.RegisterPublicRoutes(r.Group("/api/auth"), authHandler)
	auth.RegisterProtectedRoutes(r.Group("/api"), authHandler, middleware.JWTMiddleware(cfg))

	// -------------------------------------------------------------
	// Start the HTTP server on the configured port
	// Example: :8000
	// -------------------------------------------------------------
	r.Run(":" + cfg.Port)
}
