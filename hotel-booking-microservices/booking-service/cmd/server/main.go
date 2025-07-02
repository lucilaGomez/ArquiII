package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"booking-service/internal/config"
	"booking-service/internal/handlers"
	"booking-service/internal/services"
	"booking-service/pkg/mysql"
	"booking-service/pkg/memcached"
	"booking-service/pkg/amadeus"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Conectar a MySQL
	db, err := mysql.Connect(cfg.MySQLURI)
	if err != nil {
		log.Fatalf("Error conectando a MySQL: %v", err)
	}
	defer db.Close()

	// Conectar a Memcached
	mc, err := memcached.Connect(cfg.MemcachedURI)
	if err != nil {
		log.Fatalf("Error conectando a Memcached: %v", err)
	}

	// Conectar a Amadeus (corregido)
	amadeusClient := amadeus.NewClient(cfg.AmadeusBaseURL, cfg.AmadeusClientID, cfg.AmadeusClientSecret)
	
	// Verificar que Amadeus está configurado correctamente
	if cfg.AmadeusClientID == "" || cfg.AmadeusClientSecret == "" {
		log.Printf("⚠️  Warning: Credenciales de Amadeus no configuradas")
	} else {
		// Intentar obtener token para verificar conectividad
		if err := amadeusClient.GetAccessToken(); err != nil {
			log.Printf("⚠️  Warning: Error obteniendo token de Amadeus: %v", err)
		} else {
			log.Printf("✅ Cliente Amadeus configurado correctamente")
		}
	}

	// Inicializar servicios
	bookingService := services.NewBookingService(db, mc, amadeusClient, cfg.JWTSecret)

	// Inicializar handlers
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Configurar rutas
	router := setupRoutes(bookingHandler)

	// Obtener puerto
	port := cfg.Port
	log.Printf("📋 Booking Service iniciando en puerto %s", port)
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func setupRoutes(bookingHandler *handlers.BookingHandler) *gin.Engine {
	// Configurar Gin
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
			"service": "booking-service",
			"version": "1.0.0",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Rutas públicas (sin autenticación)
		api.POST("/auth/register", bookingHandler.Register)
		api.POST("/auth/login", bookingHandler.Login)
		
		// Rutas de disponibilidad (públicas)
		api.GET("/availability/:hotelId", bookingHandler.CheckAvailability)

		// Rutas protegidas (requieren autenticación)
		protected := api.Group("")
		protected.Use(bookingHandler.AuthMiddleware())
		{
			// Perfil de usuario
			protected.GET("/profile", bookingHandler.GetProfile)
			
			// Reservas
			bookings := protected.Group("/bookings")
			{
				bookings.POST("", bookingHandler.CreateBooking)                    // Crear reserva
				bookings.GET("/my-bookings", bookingHandler.GetBookings)           // NUEVA RUTA - Mis reservas
				bookings.GET("/:id", bookingHandler.GetBookingByID)                // Obtener reserva por ID
			}
		}
	}

	return router
}