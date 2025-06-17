package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Hotel estructura para los hoteles del hotel service
type Hotel struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	City        string   `json:"city"`
	Address     string   `json:"address"`
	Photos      []string `json:"photos"`
	Thumbnail   string   `json:"thumbnail"`
	Amenities   []string `json:"amenities"`
	Rating      float64  `json:"rating"`
	PriceRange  struct {
		MinPrice int    `json:"min_price"`
		MaxPrice int    `json:"max_price"`
		Currency string `json:"currency"`
	} `json:"price_range"`
	Contact struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Website string `json:"website"`
	} `json:"contact"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsActive  bool   `json:"is_active"`
}

// HotelResponse estructura para la respuesta del hotel service
type HotelResponse struct {
	Count int     `json:"count"`
	Data  []Hotel `json:"data"`
}

func main() {
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "search-service",
			"version": "1.0.0 - conectado a hotel service",
		})
	})

	router.GET("/api/v1/search/hotels", func(c *gin.Context) {
		city := c.Query("city")

		resp, err := http.Get("http://hotel_service:8080/api/v1/hotels")
		if err != nil {
			log.Printf("Error conectando con hotel service: %v", err)
			sendDemoResponse(c, city)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error leyendo respuesta del hotel service: %v", err)
			sendDemoResponse(c, city)
			return
		}

		var hotelResponse HotelResponse
		if err := json.Unmarshal(body, &hotelResponse); err != nil {
			log.Printf("Error parseando respuesta del hotel service: %v", err)
			sendDemoResponse(c, city)
			return
		}

		var filteredHotels []gin.H
		for _, hotel := range hotelResponse.Data {
			// Búsqueda mejorada que ignora acentos y mayúsculas
			if city == "" || containsIgnoreAccents(hotel.City, city) {
				filteredHotel := gin.H{
					"id":          hotel.ID,
					"name":        hotel.Name,
					"description": hotel.Description,
					"city":        strings.ToLower(hotel.City),
					"available":   true,
					"rating":      hotel.Rating,
					"min_price":   hotel.PriceRange.MinPrice,
					"max_price":   hotel.PriceRange.MaxPrice,
					"currency":    hotel.PriceRange.Currency,
					"thumbnail":   hotel.Thumbnail,
					"amenities":   hotel.Amenities,
				}
				filteredHotels = append(filteredHotels, filteredHotel)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Búsqueda exitosa - conectado a hotel service",
			"data": gin.H{
				"hotels":      filteredHotels,
				"total":       len(filteredHotels),
				"page":        1,
				"page_size":   20,
				"total_pages": 1,
			},
			"search_params": gin.H{
				"city":     city,
				"checkin":  c.Query("checkin"),
				"checkout": c.Query("checkout"),
				"guests":   c.Query("guests"),
			},
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🔍 Search Service (conectado a hotel service) iniciando en puerto %s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func sendDemoResponse(c *gin.Context, city string) {
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro 'city' es obligatorio",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Search service funcionando - búsqueda básica (modo fallback)",
		"data": gin.H{
			"hotels": []gin.H{
				{
					"id":          "demo-1",
					"name":        "Hotel Demo " + city,
					"description": "Hotel de demostración en " + city,
					"city":        city,
					"available":   true,
				},
			},
			"total":       1,
			"page":        1,
			"page_size":   20,
			"total_pages": 1,
		},
		"search_params": gin.H{
			"city": city,
		},
	})
}

// containsIgnoreAccents compara ciudades ignorando acentos y mayúsculas
func containsIgnoreAccents(haystack, needle string) bool {
	haystack = strings.ToLower(removeAccents(haystack))
	needle = strings.ToLower(removeAccents(needle))
	return strings.Contains(haystack, needle)
}

// removeAccents quita los acentos de una cadena
func removeAccents(s string) string {
	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U",
		"ñ", "n", "Ñ", "N",
	)
	return replacer.Replace(s)
}