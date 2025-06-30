package config

import (
	"os"
)

// Config contiene la configuración de la aplicación
type Config struct {
	MySQLURI           string
	MemcachedURI       string
	AmadeusClientID    string
	AmadeusClientSecret string
	AmadeusBaseURL     string
	Port               string
	Environment        string
	JWTSecret          string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	return &Config{
		MySQLURI:           getEnv("MYSQL_URI", "booking_user:booking_pass@tcp(localhost:3306)/booking_db"),
		MemcachedURI:       getEnv("MEMCACHED_URI", "localhost:11211"),
		AmadeusClientID:    getEnv("AMADEUS_CLIENT_ID", ""),
		AmadeusClientSecret: getEnv("AMADEUS_CLIENT_SECRET", ""),
		AmadeusBaseURL:     getEnv("AMADEUS_BASE_URL", "https://test.api.amadeus.com"),
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		JWTSecret:         getEnv("JWT_SECRET", "mi-secreto-super-seguro-2024"),
	}
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}