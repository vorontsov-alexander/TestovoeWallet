package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func Load() Config {
	_ = godotenv.Load("config.env")

	return Config{
		Port:   getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5433"),
		DBUser: getEnv("DB_USER", ""),
		DBPass: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	log.Println(key)
	if val, ok := os.LookupEnv(key); ok {
		log.Println(val)
		return val
	}
	return fallback
}
