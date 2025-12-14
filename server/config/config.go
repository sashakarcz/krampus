package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// OIDC Configuration
	OIDCProviderURL  string
	OIDCClientID     string
	OIDCClientSecret string
	OIDCRedirectURL  string
	OIDCScopes       []string

	// JWT Configuration
	JWTSecret string
	JWTExpiry time.Duration

	// Application Configuration
	VoteThreshold int
	AdminEmails   []string
	SyncBaseURL   string
	ServerPort    string

	// Database Configuration
	DatabasePath string
}

var AppConfig *Config

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists (silently ignore if not found)
	_ = godotenv.Load()

	config := &Config{
		// OIDC
		OIDCProviderURL:  getEnv("OIDC_PROVIDER_URL", ""),
		OIDCClientID:     getEnv("OIDC_CLIENT_ID", ""),
		OIDCClientSecret: getEnv("OIDC_CLIENT_SECRET", ""),
		OIDCRedirectURL:  getEnv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		OIDCScopes:       strings.Split(getEnv("OIDC_SCOPES", "openid,profile,email"), ","),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),

		// Application
		VoteThreshold: parseInt(getEnv("VOTE_THRESHOLD", "3")),
		AdminEmails:   parseList(getEnv("ADMIN_EMAILS", "")),
		SyncBaseURL:   getEnv("SYNC_BASE_URL", "http://localhost:8080"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),

		// Database
		DatabasePath: getEnv("DATABASE_PATH", "./database/krampus.db"),
	}

	// Validation
	if config.OIDCProviderURL == "" {
		log.Println("WARNING: OIDC_PROVIDER_URL not set - OIDC authentication will not work")
	}
	if config.OIDCClientID == "" {
		log.Println("WARNING: OIDC_CLIENT_ID not set - OIDC authentication will not work")
	}
	if config.OIDCClientSecret == "" {
		log.Println("WARNING: OIDC_CLIENT_SECRET not set - OIDC authentication will not work")
	}
	if config.JWTSecret == "change-me-in-production" {
		log.Println("WARNING: Using default JWT secret - change JWT_SECRET in production!")
	}

	AppConfig = config
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Failed to parse int from '%s', using 0\n", value)
		return 0
	}
	return i
}

func parseDuration(value string) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Failed to parse duration from '%s', using 24h\n", value)
		return 24 * time.Hour
	}
	return d
}

func parseList(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// IsAdminEmail checks if an email is in the admin list
func (c *Config) IsAdminEmail(email string) bool {
	for _, adminEmail := range c.AdminEmails {
		if strings.EqualFold(adminEmail, email) {
			return true
		}
	}
	return false
}
