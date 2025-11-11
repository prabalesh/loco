package config

import (
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
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
	AccessTokenMaxAge      int
	RefreshTokenMaxAge     int
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

		instance = &Config{
			Server: ServerConfig{
				Port: getEnv("PORT", "8080"),
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
				RefreshTokenExpiration: parseDuration("REFRESH_TOKEN_EXPIRATION", "7d"),
				AccessTokenMaxAge:      15,    // 15 minutes
				RefreshTokenMaxAge:     10080, // 7 days in minutes
			},
			Log: LogConfig{
				Level: getEnv("LOG_LEVEL", "info"),
			},
		}

		log.Println("Configuration loaded successfully")
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

func parseDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration for %s: %v", key, err)
	}
	return duration
}
