package config

import "os"

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	Port string

	// Database
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSLMode string

	// Redis
	RedisAddr string
	RedisPass string

	// JWT
	JWTSecret string
	JWTExpiry string

	// External Services
	GoogleClientID      string
	MidtransServerKey   string
	MidtransClientKey   string
	MidtransEnv         string // sandbox | production
	FCMServerKey        string
}

// Load reads config from environment with sensible defaults for development.
func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "sn_admin"),
		DBPass:            getEnv("DB_PASS", "sn_password_2026"),
		DBName:            getEnv("DB_NAME", "satria_nusantara"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:         getEnv("REDIS_PASS", ""),
		JWTSecret:         getEnv("JWT_SECRET", "sn-secret-change-in-production"),
		JWTExpiry:         getEnv("JWT_EXPIRY", "24h"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransEnv:       getEnv("MIDTRANS_ENV", "sandbox"),
		FCMServerKey:      getEnv("FCM_SERVER_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
