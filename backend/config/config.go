package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
	Version     string
	SentryDSN   string

	Database DatabaseConfig
	Worker   WorkerConfig
	Azure    AzureConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type WorkerConfig struct {
	Concurrency int
}

type AzureConfig struct {
	StorageAccount string
	StorageKey     string
	ContainerName  string
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/genshin_quiz?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: getEnv("ENVIRONMENT", "development"),
		Version:     getEnv("VERSION", "dev"),
		SentryDSN:   getEnv("SENTRY_DSN", ""),

		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "user"),
			Password:        getEnv("DB_PASSWORD", "password"),
			Name:            getEnv("DB_NAME", "genshin_quiz"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "5m"),
		},

		Worker: WorkerConfig{
			Concurrency: getEnvAsInt("WORKER_CONCURRENCY", 10),
		},

		Azure: AzureConfig{
			StorageAccount: getEnv("AZURE_STORAGE_ACCOUNT", ""),
			StorageKey:     getEnv("AZURE_STORAGE_KEY", ""),
			ContainerName:  getEnv("AZURE_CONTAINER_NAME", "genshin-quiz"),
		},

		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", "30s"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := getEnv(key, defaultValue)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	if value, err := time.ParseDuration(defaultValue); err == nil {
		return value
	}
	return time.Minute * 5
}