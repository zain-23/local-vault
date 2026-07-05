package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	// MongoDB connection
	MongoURI string
	MongoDB  string

	// JWT signing — secret signs tokens, expiry controls how long they last
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	// Resend email service
	ResendAPIKey string
	FromEmail    string

	// OAuth credentials from Google/GitHub developer console
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURL  string

	FrontendURL        string
	CORSAllowedOrigins string
}

// Load reads env vars and returns Config — call once in main()
func Load() Config {
	godotenv.Load("server/.env")

	return Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGODB_DATABASE", "localvault"),

		JWTSecret:        getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTAccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
		JWTRefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "720h")),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		FromEmail:    getEnv("FROM_EMAIL", "onboarding@resend.dev"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth/google/callback"),

		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth/github/callback"),

		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	// time.ParseDuration understands "s", "m", "h" suffixes
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute // safe default if parsing fails
	}
	return d
}