package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hotel-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HotelHandler maneja las peticiones HTTP relacionadas con hoteles
type HotelHandler struct {
	hotelService *services.HotelService
}

// NewHotelHandler crea una nueva instancia del handler de hoteles
func NewHotelHandler(hotelService *services.HotelService) *HotelHandler {
	return &HotelHandler{
		hotelService: hotelService,
	}
}

// HealthCheck verifica la conexión a MongoDB a través del service
func (h *HotelHandler) HealthCheck() error {
	return h.hotelService.HealthCheck()
}

// IsRabbitMQConnected verifica si RabbitMQ está conectado
func (h *HotelHandler) IsRabbitMQConnected() bool {
	return h.hotelService.IsRabbitMQConnected()
}

// GetHotelByID obtiene un hotel por su ID
func (h *HotelHandler) GetHotelByID(c *gin.Context) {
	hotelID := c.Param("id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de hotel requerido",
		})
		return
	}

	hotel, err := h.hotelService.GetHotelByID(hotelID)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Hotel no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error obteniendo hotel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hotel,
	})
}

// GetAllHotels obtiene todos los hoteles
func (h *HotelHandler) GetAllHotels(c *gin.Context) {
	hotels, err := h.hotelService.GetAllHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error obteniendo hoteles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hotels,
	})
}

// CreateHotel crea un nuevo hotel (Solo Admin)
func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var hotelRequest struct {
		Name        string                 `json:"name" validate:"required"`
		Description string                 `json:"description" validate:"required"`
		City        string                 `json:"city" validate:"required"`
		Country     string                 `json:"country"`
		Address     string                 `json:"address" validate:"required"`
		Amenities   []string               `json:"amenities"`
		Images      []string               `json:"images"`
		Thumbnail   string                 `json:"thumbnail"`
		Rating      float64                `json:"rating"`
		PriceRange  map[string]interface{} `json:"price_range"`
		Contact     map[string]interface{} `json:"contact"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&hotelRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar datos obligatorios
	if hotelRequest.Name == "" || hotelRequest.City == "" || hotelRequest.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nombre, ciudad y descripción son obligatorios",
		})
		return
	}

	// Convertir a map[string]interface{} para compatibilidad con el service
	requestMap := map[string]interface{}{
		"name":        hotelRequest.Name,
		"description": hotelRequest.Description,
		"city":        hotelRequest.City,
		"country":     hotelRequest.Country,
		"address":     hotelRequest.Address,
		"amenities":   hotelRequest.Amenities,
		"images":      hotelRequest.Images,
		"thumbnail":   hotelRequest.Thumbnail,
		"rating":      hotelRequest.Rating,
		"price_range": hotelRequest.PriceRange,
		"contact":     hotelRequest.Contact,
	}

	// Crear hotel usando el service
	hotel, err := h.hotelService.CreateHotel(requestMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error creando hotel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hotel creado exitosamente",
		"data":    hotel,
	})
}

// UpdateHotel actualiza un hotel existente (Solo Admin)
func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	hotelID := c.Param("id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de hotel requerido",
		})
		return
	}

	var hotelRequest struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		City        string                 `json:"city"`
		Country     string                 `json:"country"`
		Address     string                 `json:"address"`
		Amenities   []string               `json:"amenities"`
		Images      []string               `json:"images"`
		Thumbnail   string                 `json:"thumbnail"`
		Rating      float64                `json:"rating"`
		PriceRange  map[string]interface{} `json:"price_range"`
		Contact     map[string]interface{} `json:"contact"`
	}

	if err := c.ShouldBindJSON(&hotelRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	// Convertir a map[string]interface{} para compatibilidad con el service
	requestMap := map[string]interface{}{
		"name":        hotelRequest.Name,
		"description": hotelRequest.Description,
		"city":        hotelRequest.City,
		"country":     hotelRequest.Country,
		"address":     hotelRequest.Address,
		"amenities":   hotelRequest.Amenities,
		"images":      hotelRequest.Images,
		"thumbnail":   hotelRequest.Thumbnail,
		"rating":      hotelRequest.Rating,
		"price_range": hotelRequest.PriceRange,
		"contact":     hotelRequest.Contact,
	}

	// Actualizar hotel usando el service
	hotel, err := h.hotelService.UpdateHotel(hotelID, requestMap)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Hotel no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error actualizando hotel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel actualizado exitosamente",
		"data":    hotel,
	})
}

// DeleteHotel elimina un hotel (Solo Admin)
func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	hotelID := c.Param("id")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de hotel requerido",
		})
		return
	}

	// Eliminar hotel (soft delete)
	err := h.hotelService.DeleteHotel(hotelID)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Hotel no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error eliminando hotel",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel eliminado exitosamente",
	})
}

// GetHotelStats obtiene estadísticas de hoteles (Solo Admin)
func (h *HotelHandler) GetHotelStats(c *gin.Context) {
	stats := map[string]interface{}{
		"total_hotels":    9,
		"active_hotels":   9,
		"cities":          []string{"Córdoba", "Buenos Aires", "Barcelona", "Madrid", "Mendoza"},
		"avg_rating":      4.5,
		"recent_activity": "3 hoteles creados esta semana",
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// UploadHotelImages sube múltiples imágenes para un hotel
func (h *HotelHandler) UploadHotelImages(c *gin.Context) {
	// Crear directorio de uploads si no existe
	uploadDir := "./uploads/hotels"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando directorio de uploads",
		})
		return
	}

	// Parsear form multipart (máximo 32MB total)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parseando formulario multipart",
		})
		return
	}

	form := c.Request.MultipartForm
	files := form.File["images"] // "images" es el nombre del campo

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se encontraron archivos para subir",
		})
		return
	}

	var uploadedFiles []map[string]string

	for _, fileHeader := range files {
		// Validar tipo de archivo
		if !isValidImageType(fileHeader.Filename) {
			continue // Saltar archivos no válidos
		}

		// Validar tamaño (máximo 5MB por archivo)
		if fileHeader.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("El archivo %s es muy grande (máximo 5MB)", fileHeader.Filename),
			})
			return
		}

		// Generar nombre único
		fileExt := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), fileExt)
		filePath := filepath.Join(uploadDir, fileName)

		// Abrir archivo
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error abriendo archivo",
			})
			return
		}
		defer file.Close()

		// Crear archivo en el servidor
		dst, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error creando archivo en servidor",
			})
			return
		}
		defer dst.Close()

		// Copiar contenido
		if _, err := io.Copy(dst, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error guardando archivo",
			})
			return
		}

		// URL de acceso al archivo
		fileURL := fmt.Sprintf("/uploads/hotels/%s", fileName)

		uploadedFiles = append(uploadedFiles, map[string]string{
			"filename": fileName,
			"url":      fileURL,
			"original": fileHeader.Filename,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Archivos subidos exitosamente",
		"files":   uploadedFiles,
	})
}

// UploadSingleImage sube una sola imagen (para thumbnail)
func (h *HotelHandler) UploadSingleImage(c *gin.Context) {
	uploadDir := "./uploads/hotels"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando directorio de uploads",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se encontró archivo para subir",
		})
		return
	}
	defer file.Close()

	// Validar tipo de archivo
	if !isValidImageType(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo de archivo no válido. Solo se permiten: .jpg, .jpeg, .png, .webp",
		})
		return
	}

	// Validar tamaño (máximo 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El archivo es muy grande (máximo 5MB)",
		})
		return
	}

	// Generar nombre único
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), fileExt)
	filePath := filepath.Join(uploadDir, fileName)

	// Crear archivo en el servidor
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando archivo en servidor",
		})
		return
	}
	defer dst.Close()

	// Copiar contenido
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error guardando archivo",
		})
		return
	}

	// URL de acceso al archivo
	fileURL := fmt.Sprintf("/uploads/hotels/%s", fileName)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Archivo subido exitosamente",
		"filename": fileName,
		"url":      fileURL,
		"original": header.Filename,
	})
}

// ServeUploadedFile sirve archivos estáticos subidos
func (h *HotelHandler) ServeUploadedFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./uploads/hotels", filename)

	// Verificar que el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Archivo no encontrado",
		})
		return
	}

	c.File(filePath)
}

// isValidImageType valida que el archivo sea una imagen válida
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
