package config

import (
	"os"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	MongoDBURI   string
	RabbitMQURI  string
	Port         string
	Environment  string
	DatabaseName string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	return &Config{
		MongoDBURI:   getEnv("MONGODB_URI", "mongodb://admin:password123@localhost:27017/hotels_db?authSource=admin"),
		RabbitMQURI:  getEnv("RABBITMQ_URI", "amqp://admin:admin123@localhost:5672/"),
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		DatabaseName: getEnv("DATABASE_NAME", "hotels_db"),
	}
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}