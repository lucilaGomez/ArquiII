package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"hotel-service/internal/models"
	"hotel-service/internal/services"
)

// HotelHandler maneja las peticiones HTTP para hoteles
type HotelHandler struct {
	service   *services.HotelService
	validator *validator.Validate
}

// NewHotelHandler crea una nueva instancia del handler
func NewHotelHandler(service *services.HotelService) *HotelHandler {
	return &HotelHandler{
		service:   service,
		validator: validator.New(),
	}
}

// GetHotelByID obtiene un hotel por ID
func (h *HotelHandler) GetHotelByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de hotel es requerido",
		})
		return
	}

	hotel, err := h.service.GetHotelByID(id)
	if err != nil {
		if err.Error() == "hotel no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Hotel no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hotel,
	})
}

// CreateHotel crea un nuevo hotel
func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var req models.CreateHotelRequest

	// Bind JSON al struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar datos
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de validación fallidos",
			"details": err.Error(),
		})
		return
	}

	// Crear hotel
	hotel, err := h.service.CreateHotel(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando hotel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hotel creado exitosamente",
		"data": hotel,
	})
}

// UpdateHotel actualiza un hotel existente
func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de hotel es requerido",
		})
		return
	}

	var req models.UpdateHotelRequest

	// Bind JSON al struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar datos
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de validación fallidos",
			"details": err.Error(),
		})
		return
	}

	// Actualizar hotel
	hotel, err := h.service.UpdateHotel(id, &req)
	if err != nil {
		if err.Error() == "hotel no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Hotel no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error actualizando hotel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel actualizado exitosamente",
		"data": hotel,
	})
}

// GetAllHotels obtiene todos los hoteles (para testing)
func (h *HotelHandler) GetAllHotels(c *gin.Context) {
	hotels, err := h.service.GetAllHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo hoteles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hotels,
		"count": len(hotels),
	})
}