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

	// ✅ HEALTHCHECK MEJORADO - Verifica conexiones reales
	router.GET("/health", func(c *gin.Context) {
		// Verificar MongoDB
		mongoStatus := "connected"
		if err := hotelHandler.HealthCheck(); err != nil {
			mongoStatus = "disconnected"
		}

		// Verificar RabbitMQ
		rabbitStatus := "connected"
		if !hotelHandler.IsRabbitMQConnected() {
			rabbitStatus = "disconnected"
		}

		// Determinar status general
		overallStatus := "ok"
		statusCode := http.StatusOK
		if mongoStatus == "disconnected" || rabbitStatus == "disconnected" {
			overallStatus = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, gin.H{
			"status":  overallStatus,
			"service": "hotel-service",
			"version": "1.0.0",
			"checks": gin.H{
				"mongodb":  mongoStatus,
				"rabbitmq": rabbitStatus,
			},
		})
	})

	// IMPORTANTE: Servir archivos estáticos
	router.Static("/uploads", "./uploads")

	// API v1
	v1 := router.Group("/api/v1")
	{
		hotels := v1.Group("/hotels")
		{
			// Rutas públicas
			hotels.GET("/:id", hotelHandler.GetHotelByID)
			hotels.GET("", hotelHandler.GetAllHotels)
			
			// Rutas de administración (deberías agregar middleware de auth aquí)
			hotels.POST("", hotelHandler.CreateHotel)           // Crear hotel
			hotels.PUT("/:id", hotelHandler.UpdateHotel)        // Actualizar hotel
			hotels.DELETE("/:id", hotelHandler.DeleteHotel)     // Eliminar hotel
			hotels.GET("/stats", hotelHandler.GetHotelStats)    // Estadísticas
			
			// NUEVAS RUTAS para upload de imágenes
			hotels.POST("/upload-single", hotelHandler.UploadSingleImage)    // Subir imagen individual
			hotels.POST("/upload-images", hotelHandler.UploadHotelImages)    // Subir múltiples imágenes
		}
	}

	return router
}