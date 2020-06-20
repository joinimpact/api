package config

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
	DevMode          bool   // enables/disables certain debugging features for development
	JWTSecret        string // secret for signing JWTs for authorization
	Port             int    // the HTTP port to serve the application on
	DatabaseHost     string // the host of the postgres database
	DatabasePort     int    // the port of the postgres database
	DatabaseUser     string // the user of the postgres database
	DatabaseName     string // the dbname of the postgres database
	DatabasePassword string // the password to the postgres database
	SendGridAPIKey   string // the api key to access sendgrid
}

// NewConfig generates a new config from environment variables and returns a Config struct.
func NewConfig() *Config {
	return &Config{
		DevMode:          envBool("IMPACT_DEV_MODE", true),
		JWTSecret:        envString("IMPACT_JWT_SECRET", defaultJWTSecret),
		Port:             envInt("PORT", defaultPort),
		DatabaseHost:     envString("IMPACT_DATABASE_HOST", "localhost"),
		DatabasePort:     envInt("IMPACT_DATABASE_PORT", 5432),
		DatabaseUser:     envString("IMPACT_DATABASE_USER", "postgres"),
		DatabaseName:     envString("IMPACT_DATABASE_NAME", "impact"),
		DatabasePassword: envString("IMPACT_DATABASE_PASSWORD", ""),
		SendGridAPIKey:   envString("IMPACT_SENDGRID_API_KEY", ""),
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

// envInt gets an environment variable by the name provided as an int64, and defaults to the second value if no environment variable is found.
func envInt(name string, fallback int) int {
	if value := os.Getenv(name); value != "" {
		// Environment variable valid, attempt to convert to int64.
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		// No error, return converted int64.
		return int(i)
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
