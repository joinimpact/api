package core

import (
	"os"
	"strconv"
)

// defaultJWTSecret is blank for security reasons so that an error will occur if no valid JWT secret is found.
const defaultJWTSecret = ""

// defaultPort is the default port that the App will listen on.
const defaultPort = 3000

// Config stores configuration values such as API keys, secrets, and more.
type Config struct {
	DevMode   bool   // enables/disables certain debugging features for development
	JWTSecret string // secret for signing JWTs for authorization
	Port      int64  // the HTTP port to serve the application on
}

// NewConfig generates a new config from environment variables and returns a Config struct.
func NewConfig() *Config {
	return &Config{
		DevMode:   envBool("IMPACT_DEV_MODE", true),
		JWTSecret: envString("IMPACT_JWT_SECRET", defaultJWTSecret),
		Port:      envInt64("PORT", defaultPort),
	}
}

// envString gets an environment variable by the name provided as a string, and defaults to the second value if no environment variable is found.
func envString(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		// Environment variable valid, return it.
		return value
	}
	// Fallback
	return fallback
}

// envInt64 gets an environment variable by the name provided as an int64, and defaults to the second value if no environment variable is found.
func envInt64(name string, fallback int64) int64 {
	if value := os.Getenv(name); value != "" {
		// Environment variable valid, attempt to convert to int64.
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		// No error, return converted int64.
		return i
	}
	// Fallback
	return fallback
}

// envBool gets an environment variable by the name provided as a boolean, and defaults to the second value if no environment variable is found.
func envBool(name string, fallback bool) bool {
	if value := os.Getenv(name); value != "" {
		// Environment variable valid, attempt to convert to boolean.
		i, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		// No error, return converted boolean.
		return i
	}
	// Fallback
	return fallback
}
