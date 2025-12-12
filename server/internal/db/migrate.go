package db

import "database/sql"

// Migrate runs all database schema migrations needed by the application.
// If tables already exist, SQLite will ignore the creation (IF NOT EXISTS).
// This function should be called once during application startup.
func Migrate(database *sql.DB) error {

	// List of SQL migration statements.
	// New tables or alterations should be added here in order.
	migrationStatements := []string{

		// ----------------------------------------------------
		// USERS TABLE
		// Stores all user accounts with hashed passwords.
		// ----------------------------------------------------
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT,
			created_at TEXT NOT NULL
		);`,

		// ----------------------------------------------------
		// REFRESH TOKENS TABLE
		// Stores hashed refresh tokens for JWT renewal.
		// Each refresh token is linked to a user.
		// 'revoked' allows token invalidation without deletion.
		// ----------------------------------------------------
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token_hash TEXT NOT NULL,
			expires_at TEXT NOT NULL,
			revoked INTEGER DEFAULT 0,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	// Execute each migration in sequence.
	for _, statement := range migrationStatements {
		if _, err := database.Exec(statement); err != nil {
			return err // Return early if any migration fails
		}
	}

	return nil // Migrations completed successfully
}
