package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"booking-service/internal/models"
	"booking-service/internal/services"
)

// BookingHandler maneja las peticiones HTTP de reservas
type BookingHandler struct {
	bookingService *services.BookingService
	validator      *validator.Validate
}

// NewBookingHandler crea una nueva instancia del handler
func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		validator:      validator.New(),
	}
}

// Register registra un nuevo usuario
func (h *BookingHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Bind JSON
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

	// Registrar usuario
	user, err := h.bookingService.RegisterUser(&req)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "El email ya está registrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error registrando usuario",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"data": user,
	})
}

// Login autentica un usuario
func (h *BookingHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind JSON
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

	// Autenticar usuario
	user, token, err := h.bookingService.LoginUser(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciales inválidas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"data": models.AuthResponse{
			User:  user,
			Token: token,
		},
	})
}

// GetProfile obtiene el perfil del usuario autenticado
func (h *BookingHandler) GetProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token requerido",
		})
		return
	}

	user, err := h.bookingService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// CheckAvailability verifica disponibilidad de un hotel
func (h *BookingHandler) CheckAvailability(c *gin.Context) {
	hotelID := c.Param("hotelId")
	if hotelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de hotel requerido",
		})
		return
	}

	// Obtener parámetros de query
	checkInStr := c.Query("checkin")
	checkOutStr := c.Query("checkout")
	guestsStr := c.Query("guests")

	if checkInStr == "" || checkOutStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fechas de check-in y check-out son requeridas",
		})
		return
	}

	// Parsear fechas
	checkIn, err := time.Parse("2006-01-02", checkInStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de fecha inválido para check-in (usar YYYY-MM-DD)",
		})
		return
	}

	checkOut, err := time.Parse("2006-01-02", checkOutStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de fecha inválido para check-out (usar YYYY-MM-DD)",
		})
		return
	}

	// Validar fechas
	if checkIn.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La fecha de check-in no puede ser en el pasado",
		})
		return
	}

	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La fecha de check-out debe ser posterior al check-in",
		})
		return
	}

	// Parsear huéspedes
	guests := 1
	if guestsStr != "" {
		if g, err := strconv.Atoi(guestsStr); err == nil && g > 0 {
			guests = g
		}
	}

	// Crear request
	req := &models.AvailabilityRequest{
		HotelID:      hotelID,
		CheckInDate:  checkIn,
		CheckOutDate: checkOut,
		Guests:       guests,
	}

	// Verificar disponibilidad
	availability, err := h.bookingService.CheckAvailability(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error verificando disponibilidad",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": availability,
	})
}

// CreateBooking crea una nueva reserva
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token requerido",
		})
		return
	}

	var req models.CreateBookingRequest

	// Bind JSON
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

	// Validar fechas
	if req.CheckInDate.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La fecha de check-in no puede ser en el pasado",
		})
		return
	}

	if req.CheckOutDate.Before(req.CheckInDate) || req.CheckOutDate.Equal(req.CheckInDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La fecha de check-out debe ser posterior al check-in",
		})
		return
	}

	// Crear reserva
	booking, err := h.bookingService.CreateBooking(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "no disponible") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creando reserva",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reserva creada exitosamente",
		"data": booking,
	})
}

// GetBookings obtiene las reservas del usuario autenticado
func (h *BookingHandler) GetBookings(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token requerido",
		})
		return
	}

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo reservas",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
		"count": len(bookings),
	})
}

// GetBookingByID obtiene una reserva específica
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token requerido",
		})
		return
	}

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de reserva inválido",
		})
		return
	}

	booking, err := h.bookingService.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Reserva no encontrada",
		})
		return
	}

	// Verificar que la reserva pertenezca al usuario
	if booking.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No tienes acceso a esta reserva",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": booking,
	})
}

// AuthMiddleware middleware de autenticación
func (h *BookingHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token de autorización requerido",
			})
			c.Abort()
			return
		}

		// Verificar formato Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token
		userID, err := h.bookingService.ValidateJWTToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido",
			})
			c.Abort()
			return
		}

		// Guardar userID en el contexto
		c.Set("userID", userID)
		c.Next()
	}
}

// getUserIDFromContext obtiene el ID del usuario del contexto
func (h *BookingHandler) getUserIDFromContext(c *gin.Context) int {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}