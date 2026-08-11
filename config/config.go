package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort       string
	AppEnv        string
	UploadDir     string
	PublicBaseURL string
	DBURL         string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiry     string
}

// load config dari env
func Load() Config {
	return Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "development"),
		UploadDir:     getEnv("UPLOAD_DIR", "uploads"),
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", "http://localhost:8080"),
		DBURL:         getEnv("DATABASE_URL", ""),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "my_api"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		JWTExpiry:     getEnv("JWT_EXPIRY", "24h"),
	}
}

// config database
func (c Config) DatabaseURL() string {
	if c.DBURL != "" {
		return c.DBURL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

// cek ada nilai env tidak jika tidak ada isi dengan nilai default
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
