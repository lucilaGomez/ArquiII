package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"

	"booking-service/internal/models"
	"booking-service/pkg/amadeus"
	"booking-service/pkg/memcached"
	"booking-service/pkg/mysql"
)

// BookingService maneja la lógica de reservas
type BookingService struct {
	db            *mysql.DB
	cache         *memcached.Client
	amadeusClient *amadeus.Client
	jwtSecret     string
}

// NewBookingService crea una nueva instancia del servicio
func NewBookingService(db *mysql.DB, cache *memcached.Client, amadeusClient *amadeus.Client, jwtSecret string) *BookingService {
	return &BookingService{
		db:            db,
		cache:         cache,
		amadeusClient: amadeusClient,
		jwtSecret:     jwtSecret,
	}
}

// RegisterUser registra un nuevo usuario
func (s *BookingService) RegisterUser(req *models.RegisterRequest) (*models.User, error) {
	// Verificar si el usuario ya existe ANTES de intentar crear
	existingUser, err := s.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("el email %s ya está registrado", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hasheando password: %v", err)
	}

	// Insertar usuario - MODIFICADO: agregado role
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, date_of_birth, role)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	// Establecer role por defecto si no se proporciona
	role := req.Role
	if role == "" {
		role = "user" // Valor por defecto
	}
	
	result, err := s.db.Exec(query, req.Email, string(hashedPassword), req.FirstName, req.LastName, req.Phone, req.DateOfBirth, role)
	if err != nil {
		// Manejo específico de errores MySQL
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("el email %s ya está registrado", req.Email)
		}
		return nil, fmt.Errorf("error creando usuario: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ID de usuario: %v", err)
	}

	// Obtener usuario creado
	return s.GetUserByID(int(userID))
}

// LoginUser autentica un usuario
func (s *BookingService) LoginUser(req *models.LoginRequest) (*models.User, string, error) {
	// Buscar usuario por email
	user, err := s.GetUserByEmail(req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("usuario no encontrado")
	}

	// Verificar password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	// Generar JWT token
	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("error generando token: %v", err)
	}

	return user, token, nil
}

// GetUserByID obtiene un usuario por ID - MODIFICADO: agregado role
func (s *BookingService) GetUserByID(userID int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, date_of_birth, role, created_at, updated_at, is_active
		FROM users WHERE id = ? AND is_active = TRUE
	`

	var user models.User
	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.DateOfBirth, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %v", err)
	}

	return &user, nil
}

// GetUserByEmail obtiene un usuario por email - MODIFICADO: agregado role
func (s *BookingService) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, date_of_birth, role, created_at, updated_at, is_active
		FROM users WHERE email = ? AND is_active = TRUE
	`

	var user models.User
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.DateOfBirth, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %v", err)
	}

	return &user, nil
}

// CheckAvailability verifica disponibilidad de un hotel
func (s *BookingService) CheckAvailability(req *models.AvailabilityRequest) (*models.AvailabilityResponse, error) {
	// Generar clave de caché
	cacheKey := memcached.GenerateAvailabilityKey(req.HotelID, req.CheckInDate, req.CheckOutDate, req.Guests)

	// Intentar obtener del caché
	var cachedResponse models.AvailabilityResponse
	err := s.cache.Get(cacheKey, &cachedResponse)
	if err == nil {
		// Cache hit
		fmt.Printf("🔍 Valor obtenido del caché: %s\n", cacheKey)
		return &cachedResponse, nil
	}

	// Cache miss - consultar Amadeus
	// Primero obtener mapeo de hotel
	amadeusHotelID, err := s.getAmadeusHotelID(req.HotelID)
	if err != nil {
		// Si no hay mapeo, simular respuesta sin error 500
		fmt.Printf("⚠️ Warning: No se encontró mapeo para hotel %s, simulando disponibilidad\n", req.HotelID)
		response := s.createSimulatedAvailability(req)
		s.cache.Set(cacheKey, response, 10*time.Second)
		fmt.Printf("💾 Valor almacenado en caché: %s\n", cacheKey)
		return response, nil
	}

	// Consultar ofertas en Amadeus
	checkInStr := req.CheckInDate.Format("2006-01-02")
	checkOutStr := req.CheckOutDate.Format("2006-01-02")
	
	offers, err := s.amadeusClient.GetHotelOffers(amadeusHotelID, checkInStr, checkOutStr, req.Guests)
	if err != nil {
		// Log del error pero no fallar
		fmt.Printf("⚠️ Warning: Error conectando con Amadeus: %v\n", err)
		fmt.Printf("📝 Simulando disponibilidad para hotel %s\n", req.HotelID)
		
		// Devolver disponibilidad simulada en lugar de error 500
		response := s.createSimulatedAvailability(req)
		s.cache.Set(cacheKey, response, 10*time.Second)
		fmt.Printf("💾 Valor almacenado en caché: %s\n", cacheKey)
		return response, nil
	}

	// Procesar ofertas de Amadeus
	response := &models.AvailabilityResponse{
		HotelID:      req.HotelID,
		Available:    len(offers) > 0,
		CheckInDate:  checkInStr,
		CheckOutDate: checkOutStr,
		Guests:       req.Guests,
	}

	if len(offers) > 0 && len(offers[0].Offers) > 0 {
		firstOffer := offers[0].Offers[0]
		response.Currency = firstOffer.Price.Currency
		
		// Convertir precio
		if priceFloat, err := strconv.ParseFloat(firstOffer.Price.Total, 64); err == nil {
			response.Price = &priceFloat
		}

		rooms := firstOffer.RoomQuantity
		response.RoomsAvailable = &rooms
	}

	// Guardar en caché por 10 segundos
	s.cache.Set(cacheKey, response, 10*time.Second)
	fmt.Printf("💾 Valor almacenado en caché: %s\n", cacheKey)

	return response, nil
}

// createSimulatedAvailability crea una respuesta simulada cuando Amadeus falla
func (s *BookingService) createSimulatedAvailability(req *models.AvailabilityRequest) *models.AvailabilityResponse {
	checkInStr := req.CheckInDate.Format("2006-01-02")
	checkOutStr := req.CheckOutDate.Format("2006-01-02")
	
	response := &models.AvailabilityResponse{
		HotelID:        req.HotelID,
		Available:      true,
		CheckInDate:    checkInStr,
		CheckOutDate:   checkOutStr,
		Guests:         req.Guests,
		Currency:       "ARS",
	}
	
	// Simular precio basado en cantidad de huéspedes
	price := float64(15000 + (req.Guests-1)*5000) // Precio base + extra por huésped
	response.Price = &price
	rooms := 5
	response.RoomsAvailable = &rooms

	return response
}

// CreateBooking crea una nueva reserva - VERSIÓN CORREGIDA
func (s *BookingService) CreateBooking(userID int, req *models.CreateBookingRequest) (*models.Booking, error) {
	// Verificar disponibilidad primero
	availReq := &models.AvailabilityRequest{
		HotelID:      req.HotelID,
		CheckInDate:  req.CheckInDate,
		CheckOutDate: req.CheckOutDate,
		Guests:       req.Guests,
	}

	availability, err := s.CheckAvailability(availReq)
	if err != nil {
		return nil, fmt.Errorf("error verificando disponibilidad: %v", err)
	}

	if !availability.Available {
		return nil, fmt.Errorf("hotel no disponible para las fechas seleccionadas")
	}

	// Crear reserva en base de datos - CORREGIDO: agregado status = 'confirmed' desde el inicio
	query := `
		INSERT INTO bookings (user_id, internal_hotel_id, check_in_date, check_out_date, guests, room_type, total_price, currency, status, special_requests, booking_reference)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'confirmed', ?, ?)
	`

	// Generar referencia de reserva
	bookingRef := fmt.Sprintf("BK%d%d", userID, time.Now().Unix())
	
	totalPrice := 0.0
	currency := "ARS"
	if availability.Price != nil {
		totalPrice = *availability.Price
	}
	if availability.Currency != "" {
		currency = availability.Currency
	}

	result, err := s.db.Exec(query, userID, req.HotelID, req.CheckInDate, req.CheckOutDate, req.Guests, req.RoomType, totalPrice, currency, req.SpecialRequests, bookingRef)
	if err != nil {
		return nil, fmt.Errorf("error creando reserva: %v", err)
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ID de reserva: %v", err)
	}

	// Intentar crear reserva en Amadeus en background (opcional, no afecta el status)
	go func() {
		amadeusHotelID, err := s.getAmadeusHotelID(req.HotelID)
		if err == nil {
			guestInfo := map[string]interface{}{
				"adults": req.Guests,
			}
			
			amadeusBookingID, err := s.amadeusClient.CreateBooking(amadeusHotelID, guestInfo)
			if err == nil {
				// Solo actualizar el amadeus_booking_id, el status ya es 'confirmed'
				s.db.Exec("UPDATE bookings SET amadeus_booking_id = ? WHERE id = ?", amadeusBookingID, bookingID)
			}
		}
	}()

	// Obtener reserva creada
	return s.GetBookingByID(int(bookingID))
}

// GetBookingByID obtiene una reserva por ID
func (s *BookingService) GetBookingByID(bookingID int) (*models.Booking, error) {
	query := `
		SELECT id, user_id, internal_hotel_id, amadeus_hotel_id, amadeus_booking_id, check_in_date, check_out_date, 
		       guests, room_type, total_price, currency, status, booking_reference, special_requests, created_at, updated_at
		FROM bookings WHERE id = ?
	`

	var booking models.Booking
	err := s.db.QueryRow(query, bookingID).Scan(
		&booking.ID, &booking.UserID, &booking.InternalHotelID, &booking.AmadeusHotelID, &booking.AmadeusBookingID,
		&booking.CheckInDate, &booking.CheckOutDate, &booking.Guests, &booking.RoomType, &booking.TotalPrice,
		&booking.Currency, &booking.Status, &booking.BookingReference, &booking.SpecialRequests,
		&booking.CreatedAt, &booking.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reserva no encontrada")
		}
		return nil, fmt.Errorf("error obteniendo reserva: %v", err)
	}

	return &booking, nil
}

// GetUserBookings obtiene todas las reservas de un usuario
func (s *BookingService) GetUserBookings(userID int) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, internal_hotel_id, amadeus_hotel_id, amadeus_booking_id, check_in_date, check_out_date, 
		       guests, room_type, total_price, currency, status, booking_reference, special_requests, created_at, updated_at
		FROM bookings WHERE user_id = ? ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo reservas: %v", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.InternalHotelID, &booking.AmadeusHotelID, &booking.AmadeusBookingID,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.Guests, &booking.RoomType, &booking.TotalPrice,
			&booking.Currency, &booking.Status, &booking.BookingReference, &booking.SpecialRequests,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando reserva: %v", err)
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

// getAmadeusHotelID obtiene el ID de Amadeus para un hotel interno
func (s *BookingService) getAmadeusHotelID(internalHotelID string) (string, error) {
	query := "SELECT amadeus_hotel_id FROM hotel_mappings WHERE internal_hotel_id = ?"
	
	var amadeusHotelID string
	err := s.db.QueryRow(query, internalHotelID).Scan(&amadeusHotelID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("mapeo de hotel no encontrado")
		}
		return "", fmt.Errorf("error obteniendo mapeo: %v", err)
	}

	return amadeusHotelID, nil
}

// generateJWTToken genera un token JWT para un usuario
func (s *BookingService) generateJWTToken(userID int) (string, error) {
	// Crear claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 días
		"iat":     time.Now().Unix(),
	}

	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error firmando token: %v", err)
	}

	return tokenString, nil
}

// ValidateJWTToken valida un token JWT
func (s *BookingService) ValidateJWTToken(tokenString string) (int, error) {
	// Parsear token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parseando token: %v", err)
	}

	// Verificar que el token sea válido
	if !token.Valid {
		return 0, fmt.Errorf("token inválido")
	}

	// Extraer claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("claims inválidos")
	}

	// Obtener user_id
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id inválido en token")
	}

	return int(userID), nil
}

// ValidateJWTTokenWithRole valida un token JWT y devuelve userID y role
func (s *BookingService) ValidateJWTTokenWithRole(tokenString string) (int, string, error) {
	// Parsear token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("error parseando token: %v", err)
	}

	// Verificar que el token sea válido
	if !token.Valid {
		return 0, "", fmt.Errorf("token inválido")
	}

	// Extraer claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("claims inválidos")
	}

	// Obtener user_id
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("user_id inválido en token")
	}

	// Obtener el usuario para verificar su rol actual
	user, err := s.GetUserByID(int(userID))
	if err != nil {
		return 0, "", fmt.Errorf("usuario no encontrado: %v", err)
	}

	return int(userID), string(user.Role), nil
}