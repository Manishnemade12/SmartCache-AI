package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	RedisURL    string
	GeminiKey   string
	WorkerCount int
	CacheTTL    int // seconds
}

var C Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	C = Config{
		Port:        getEnv("PORT", "8080"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		GeminiKey:   getEnv("GEMINI_API_KEY", ""),
		WorkerCount: getEnvInt("WORKER_COUNT", 3),
		CacheTTL:    getEnvInt("CACHE_TTL", 300),
	}

	if C.GeminiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY is not set — AI summarization will fail")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
