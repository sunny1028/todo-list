package config

import "os"

type Config struct {
	Port       string
	DBPath     string
	CorsOrigin string
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "todo.db"),
		CorsOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
