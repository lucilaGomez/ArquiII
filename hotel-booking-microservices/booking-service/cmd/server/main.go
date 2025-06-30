package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	
	"booking-service/internal/config"
	"booking-service/internal/handlers"
	"booking-service/internal/services"
	"booking-service/pkg/amadeus"
	"booking-service/pkg/memcached"
	"booking-service/pkg/mysql"
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
	cache, err := memcached.Connect(cfg.MemcachedURI)
	if err != nil {
		log.Fatalf("Error conectando a Memcached: %v", err)
	}
	defer cache.Close()

	// Inicializar cliente de Amadeus
	var amadeusClient *amadeus.Client
	if cfg.AmadeusClientID != "" && cfg.AmadeusClientSecret != "" {
		amadeusClient = amadeus.NewClient(cfg.AmadeusBaseURL, cfg.AmadeusClientID, cfg.AmadeusClientSecret)
		log.Println("🌐 Cliente de Amadeus inicializado")
	} else {
		log.Println("⚠️ Credenciales de Amadeus no configuradas - funcionando en modo simulado")
		amadeusClient = amadeus.NewClient(cfg.AmadeusBaseURL, "demo", "demo")
	}

	// Inicializar servicios
	bookingService := services.NewBookingService(db, cache, amadeusClient, cfg.JWTSecret)

	// Inicializar handlers
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Configurar rutas
	router := setupRoutes(bookingHandler)

	// Obtener puerto
	port := cfg.Port
	log.Printf("👥 Booking Service iniciando en puerto %s", port)
	
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

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Rutas públicas (auth)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", bookingHandler.Register)
			auth.POST("/login", bookingHandler.Login)
		}

		// Rutas de disponibilidad (públicas)
		availability := v1.Group("/availability")
		{
			availability.GET("/:hotelId", bookingHandler.CheckAvailability)
		}

		// Rutas protegidas (requieren autenticación)
		protected := v1.Group("/")
		protected.Use(bookingHandler.AuthMiddleware())
		{
			// Perfil de usuario
			protected.GET("/profile", bookingHandler.GetProfile)

			// Reservas
			bookings := protected.Group("/bookings")
			{
				bookings.POST("", bookingHandler.CreateBooking)
				bookings.GET("", bookingHandler.GetBookings)
				bookings.GET("/:id", bookingHandler.GetBookingByID)
			}
		}
	}

	return router
}