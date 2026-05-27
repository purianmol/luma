package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppEnv string

const (
	EnvDevelopment AppEnv = "development"
	EnvStaging     AppEnv = "staging"
	EnvProduction  AppEnv = "production"
)

type Config struct {
	// ── App ──────────────────────────────────────────────────

	// AppEnv controls log verbosity, Gin mode, and validation strictness.
	AppEnv AppEnv

	// Port is the TCP port this service listens on.
	Port string

	// LogLevel controls zap's minimum log level: "debug", "info", "warn", "error".
	LogLevel string

	// ── Database (PostgreSQL) ─────────────────────────────────

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// DBDSN is assembled from the fields above. Use this to open a GORM connection.
	// It is never set directly from an env var — it's always computed in Load().
	DBDSN string

	// ── Cache (Redis) ─────────────────────────────────────────

	RedisHost     string
	RedisPort     string
	RedisPassword string // empty string = no auth (acceptable in dev)

	// RedisAddr is "host:port", the format the Redis client expects.
	RedisAddr string

	// ── Authentication ────────────────────────────────────────

	// JWTSecret is the HMAC-SHA256 key used to sign access tokens.
	// An empty secret makes tokens trivially forgeable — always required.
	JWTSecret string

	// JWTExpiryHours is the access token lifetime. 24 hours is a sensible default:
	// short enough to limit damage from theft, long enough for a full work day.
	JWTExpiryHours int

	// RefreshTokenExpiryDays is the refresh token lifetime.
	// Refresh tokens are long-lived because they live only in HttpOnly cookies
	// and are stored in the database (where they can be revoked).
	RefreshTokenExpiryDays int

	// ── Service URLs (used by the API gateway for routing) ────

	AuthServiceURL string
	ChatServiceURL string
	FileServiceURL string

	// ── Rate limiting ─────────────────────────────────────────

	// RateLimitRPS is the maximum requests per second per IP address.
	RateLimitRPS int

	// ── File storage ──────────────────────────────────────────

	// StorageRoot is the filesystem directory where uploaded files are written.
	// In production this would be replaced by an S3 bucket name, but the
	// local filesystem works fine for development and integration tests.
	StorageRoot string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:   AppEnv(getEnv("APP_ENV", string(EnvDevelopment))),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "collab"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "secret"),
		DBName:     getEnv("POSTGRES_DB", "collabdb"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		JWTSecret:              getEnv("JWT_SECRET", ""),
		JWTExpiryHours:         getEnvInt("JWT_EXPIRY_HOURS", 24),
		RefreshTokenExpiryDays: getEnvInt("REFRESH_TOKEN_EXPIRY_DAYS", 30),

		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		ChatServiceURL: getEnv("CHAT_SERVICE_URL", "http://localhost:8082"),
		FileServiceURL: getEnv("FILE_SERVICE_URL", "http://localhost:8083"),

		RateLimitRPS: getEnvInt("RATE_LIMIT_RPS", 100),
		StorageRoot:  getEnv("STORAGE_ROOT", "/tmp/collab-files"),
	}
	cfg.DBDSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	cfg.RedisAddr = fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil

}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("config: %s must be an integer, got %q: %v", key, v, err))
	}
	return n
}
