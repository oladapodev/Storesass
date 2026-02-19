package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBDriver  string
	DBDSN     string
	RedisURL  string
	JWTSecret string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBDriver:  getEnv("DB_DRIVER", "sqlite"),
		DBDSN:     getEnv("DB_DSN", "storefront.db"),
		RedisURL:  getEnv("REDIS_URL", ""),
		JWTSecret: getEnv("JWT_SECRET", "changeme"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
