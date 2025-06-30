package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"search-service/internal/config"
	"search-service/internal/handlers"
	"search-service/internal/services"
	"search-service/pkg/rabbitmq"
	"search-service/pkg/solr"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Configurar Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Printf("🔍 Search Service iniciando...")
	log.Printf("📍 Solr URL: %s", cfg.SolrURI)
	log.Printf("🐰 RabbitMQ URL: %s", cfg.RabbitMQURI)

	// Inicializar Solr
	solrClient := solr.NewClient(cfg.SolrURI)
	log.Println("🔎 Cliente Solr inicializado")

	// Verificar conexión a Solr
	if err := solrClient.Ping(); err != nil {
		log.Printf("⚠️ No se pudo conectar a Solr: %v", err)
		log.Println("🔄 Continuando en modo fallback...")
	} else {
		log.Println("✅ Conexión a Solr exitosa")
	}

	// Inicializar RabbitMQ Consumer
	rabbitConsumer, err := rabbitmq.NewConsumer(cfg.RabbitMQURI)
	if err != nil {
		log.Printf("⚠️ No se pudo conectar a RabbitMQ: %v", err)
		log.Println("🔄 Continuando sin consumer de eventos...")
		rabbitConsumer = nil
	} else {
		log.Println("✅ RabbitMQ Consumer conectado")
	}

	// Inicializar servicio
	searchService := services.NewSearchService(solrClient, rabbitConsumer)

	// Iniciar consumer de eventos si está disponible
	if rabbitConsumer != nil {
		go func() {
			if err := searchService.StartEventConsumer(); err != nil {
				log.Printf("❌ Error en consumer de eventos: %v", err)
			}
		}()
		log.Println("🚀 Consumer de eventos iniciado")
	}

	// Inicializar handlers
	searchHandler := handlers.NewSearchHandler(searchService)

	// Configurar rutas
	router := setupRoutes(searchHandler)

	// Configurar graceful shutdown
	setupGracefulShutdown(rabbitConsumer)

	port := cfg.Port
	log.Printf("🌐 Search Service iniciando en puerto %s", port)
	log.Println("🎯 Funcionalidades activas:")
	log.Println("   - Búsqueda en Solr")
	log.Println("   - Verificación concurrente de disponibilidad")
	if rabbitConsumer != nil {
		log.Println("   - Sincronización automática vía RabbitMQ")
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

func setupRoutes(searchHandler *handlers.SearchHandler) *gin.Engine {
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

	// Middleware para logging y recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check mejorado
	router.GET("/health", searchHandler.HealthCheck)

	// API v1
	v1 := router.Group("/api/v1")
	{
		search := v1.Group("/search")
		{
			search.GET("/hotels", searchHandler.SearchHotels)
		}
	}

	return router
}

func setupGracefulShutdown(rabbitConsumer *rabbitmq.Consumer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("🛑 Cerrando Search Service...")
		
		if rabbitConsumer != nil {
			rabbitConsumer.Close()
		}
		
		log.Println("✅ Search Service cerrado correctamente")
		os.Exit(0)
	}()
}