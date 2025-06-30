package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"search-service/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// SearchHotels maneja las búsquedas de hoteles
func (h *SearchHandler) SearchHotels(c *gin.Context) {
	// Obtener parámetros de búsqueda
	city := c.Query("city")
	query := c.Query("q")
	checkin := c.Query("checkin")
	checkout := c.Query("checkout")
	guests := c.Query("guests")

	// Si no hay ningún criterio de búsqueda, devolver error
	if city == "" && query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Debe proporcionar al menos un criterio de búsqueda (city o q)",
		})
		return
	}

	// Realizar búsqueda con disponibilidad
	results, err := h.searchService.SearchHotelsWithAvailability(query, city, checkin, checkout, guests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error realizando búsqueda",
			"details": err.Error(),
		})
		return
	}

	// Respuesta exitosa
	response := gin.H{
		"success": true,
		"message": "Búsqueda exitosa con Solr y disponibilidad concurrente",
		"data":    results,
		"search_params": gin.H{
			"city":     city,
			"query":    query,
			"checkin":  checkin,
			"checkout": checkout,
			"guests":   guests,
		},
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck endpoint para verificar estado del servicio
func (h *SearchHandler) HealthCheck(c *gin.Context) {
	health := h.searchService.HealthCheck()
	
	// Determinar status code basado en la salud de los componentes
	statusCode := http.StatusOK
	for _, status := range health {
		if status == "down" {
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	c.JSON(statusCode, gin.H{
		"status":     "ok",
		"service":    "search-service",
		"version":    "1.0.0",
		"components": health,
		"features": []string{
			"Solr search engine",
			"RabbitMQ event consumer",
			"Concurrent availability checking",
			"Hotel data synchronization",
		},
	})
}