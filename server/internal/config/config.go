package config

import (
	"os"
	"strconv"
)

// Config holds all environment-driven configuration required
// for running the application. These values are loaded once at startup.
type Config struct {
	Port                  string // Port the HTTP server listens on (e.g., "8000")
	JWTSecret             string // Secret used to sign JWT access tokens
	DBPath                string // Path to SQLite database file
	DataDir               string // Directory where app data is stored (e.g., SQLite file)
	CookieDomain          string // Domain for setting cookies (e.g., "localhost")
	CookieSecure          bool   // Whether cookies require HTTPS (true in production)
	AccessExpiry          int    // Access token lifetime in seconds
	RefreshExpiry         int    // Refresh token lifetime in seconds
	EncryptionKey         string // Server-side encryption key for notes
	AppBaseURL            string // Base URL of the frontend app
	EncryptedNotesEnabled bool
	UserSaltLength        int
}

// getString retrieves a string value from the environment.
// If the environment variable is missing or empty, it returns the fallback default.
func getString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getInt retrieves an integer from the environment.
// If missing or invalid, it returns the fallback default.
func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

// LoadFromEnv reads all configuration values from environment variables
// and returns a fully initialized Config struct.
// Default values are used when variables are not provided.
func LoadFromEnv() *Config {
	return &Config{
		Port:                  getString("PORT", "8000"),
		JWTSecret:             getString("JWT_SECRET", "dev_secret"),
		DBPath:                getString("DB_PATH", "./data/app.db"),
		DataDir:               getString("DATA_DIR", "./data"),
		CookieDomain:          getString("COOKIE_DOMAIN", "localhost"),
		CookieSecure:          getString("COOKIE_SECURE", "false") == "true",
		EncryptionKey:         getString("ENCRYPTION_KEY", ""),
		AppBaseURL:            getString("AppBaseURL", ""),
		EncryptedNotesEnabled: getString("ENCRYPTED_NOTES_ENABLED", "true") == "true",
		AccessExpiry:          getInt("ACCESS_EXPIRY", 300),
		RefreshExpiry:         getInt("REFRESH_EXPIRY", 604800),
		UserSaltLength:        getInt("ENCRYPTION_USER_SALT_LENGTH", 16),
	}
}
