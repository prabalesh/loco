package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	CORS     CORSConfig
	Email    EmailConfig
	Log      LogConfig
}

type CORSConfig struct {
	AllowedOrigins []string
}

type ServerConfig struct {
	Port string
	Env  string // "development" or "production"
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	AccessTokenMaxAge      int // in seconds
	RefreshTokenMaxAge     int // in seconds
}

type CookieConfig struct {
	Secure   bool
	SameSite string
	Domain   string
}

type EmailConfig struct {
	ResendAPIKey          string
	FromEmail             string
	FromName              string
	OTPExpirationMinutes  int
	MaxOTPAttempts        int
	ResendCooldownMinutes int
}

type LogConfig struct {
	Level string
}

var (
	instance *Config
	once     sync.Once
)

// InitConfig initializes the singleton config instance
func InitConfig() {
	once.Do(func() {
		log.Println("Initializing configuration...")

		env := getEnv("ENV", "development")
		isProduction := env == "production"

		allowedOrigins := parseAllowedOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"))

		instance = &Config{
			Server: ServerConfig{
				Port: getEnv("PORT", "8080"),
				Env:  env,
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: mustGetEnv("DB_PASSWORD"),
				Name:     getEnv("DB_NAME", "coding_platform"),
				SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			},
			JWT: JWTConfig{
				AccessTokenSecret:      mustGetEnv("ACCESS_TOKEN_SECRET"),
				RefreshTokenSecret:     mustGetEnv("REFRESH_TOKEN_SECRET"),
				AccessTokenExpiration:  parseDuration("ACCESS_TOKEN_EXPIRATION", "15m"),
				RefreshTokenExpiration: parseDuration("REFRESH_TOKEN_EXPIRATION", "168h"), // 7 days
				AccessTokenMaxAge:      900,                                               // 15 minutes in seconds
				RefreshTokenMaxAge:     604800,                                            // 7 days in seconds
			},
			Cookie: CookieConfig{
				Secure:   isProduction,
				SameSite: getEnv("COOKIE_SAMESITE", getDefaultSameSite(isProduction)),
				Domain:   getEnv("COOKIE_DOMAIN", ""),
			},
			CORS: CORSConfig{
				AllowedOrigins: allowedOrigins,
			},
			Email: EmailConfig{
				ResendAPIKey:          mustGetEnv("RESEND_API_KEY"),
				FromEmail:             getEnv("FROM_EMAIL", "onboarding@yourdomain.com"),
				FromName:              getEnv("FROM_NAME", "Loco Platform"),
				OTPExpirationMinutes:  parseInt("OTP_EXPIRATION_MINUTES", 10),
				MaxOTPAttempts:        parseInt("MAX_OTP_ATTEMPTS", 5),
				ResendCooldownMinutes: parseInt("RESEND_COOLDOWN_MINUTES", 2),
			},
			Log: LogConfig{
				Level: getEnv("LOG_LEVEL", "info"),
			},
		}

		log.Println("Configuration loaded successfully")
		log.Printf("Environment: %s", env)
		log.Printf("Cookie Secure: %v", instance.Cookie.Secure)
	})
}

// GetConfig returns the singleton config instance
func GetConfig() *Config {
	if instance == nil {
		panic("config not initialized, call InitConfig() first")
	}
	return instance
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// mustGetEnv gets required environment variable or panics
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// parseDuration parses duration from environment variable
func parseDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration for %s: %v", key, err)
	}
	return duration
}

// parseInt parses integer from environment variable
func parseInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid integer for %s: %v", key, err)
	}
	return intValue
}

// getDefaultSameSite returns default SameSite policy based on environment
func getDefaultSameSite(isProduction bool) string {
	if isProduction {
		return "strict"
	}
	return "lax"
}

func parseAllowedOrigins(originsStr string) []string {
	origins := strings.Split(originsStr, ",")
	result := make([]string, 0, len(origins))

	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
