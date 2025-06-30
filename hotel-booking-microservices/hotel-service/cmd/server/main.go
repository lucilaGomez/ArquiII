package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"hotel-service/internal/config"
	"hotel-service/internal/handlers"
	"hotel-service/internal/services"
	"hotel-service/pkg/mongodb"
	"hotel-service/pkg/rabbitmq"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Conectar a MongoDB
	mongoClient, err := mongodb.Connect(cfg.MongoDBURI)
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer mongodb.Disconnect(mongoClient)

	// Conectar a RabbitMQ
	rabbitConn, err := rabbitmq.Connect(cfg.RabbitMQURI)
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Inicializar servicios
	hotelService := services.NewHotelService(mongoClient, rabbitConn)

	// Inicializar handlers
	hotelHandler := handlers.NewHotelHandler(hotelService)

	// Configurar rutas
	router := setupRoutes(hotelHandler)

	// Obtener puerto
	port := cfg.Port
	log.Printf("🏨 Hotel Service iniciando en puerto %s", port)
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func setupRoutes(hotelHandler *handlers.HotelHandler) *gin.Engine {
	// Configurar Gin en modo release para producción
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware para CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Middleware para logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Ruta de salud
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "hotel-service",
			"version": "1.0.0",
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		hotels := v1.Group("/hotels")
		{
			hotels.GET("/:id", hotelHandler.GetHotelByID)
			hotels.POST("", hotelHandler.CreateHotel)
			hotels.PUT("/:id", hotelHandler.UpdateHotel)
			hotels.GET("", hotelHandler.GetAllHotels)
		}
	}

	return router
}