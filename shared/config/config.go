package config

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
