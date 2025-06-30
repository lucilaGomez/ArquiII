package config

import (
	"os"
)

// Config contiene la configuración de la aplicación
type Config struct {
	SolrURI         string
	RabbitMQURI     string
	HotelServiceURI string
	Port            string
	Environment     string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	return &Config{
		SolrURI:         getEnv("SOLR_URI", "http://localhost:8983/solr/hotels"),
		RabbitMQURI:     getEnv("RABBITMQ_URI", "amqp://admin:admin123@localhost:5672/"),
		HotelServiceURI: getEnv("HOTEL_SERVICE_URI", "http://localhost:8001"),
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}