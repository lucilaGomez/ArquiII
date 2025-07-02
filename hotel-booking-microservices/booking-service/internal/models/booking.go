package models

import (
	"time"
)

// UserRole define los roles de usuario  
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User representa un usuario del sistema
type User struct {
	ID          int       `json:"id" db:"id"`
	Email       string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string   `json:"-" db:"password_hash"`
	FirstName   string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName    string    `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Role        UserRole  `json:"role" db:"role"` // ← LÍNEA AGREGADA
	Phone       string    `json:"phone" db:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

// IsAdmin verifica si el usuario es administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser verifica si el usuario es usuario normal
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}

// Booking representa una reserva
type Booking struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	InternalHotelID  string    `json:"internal_hotel_id" db:"internal_hotel_id"`
	AmadeusHotelID   *string   `json:"amadeus_hotel_id" db:"amadeus_hotel_id"`
	AmadeusBookingID *string   `json:"amadeus_booking_id" db:"amadeus_booking_id"`
	CheckInDate      time.Time `json:"check_in_date" db:"check_in_date"`
	CheckOutDate     time.Time `json:"check_out_date" db:"check_out_date"`
	Guests           int       `json:"guests" db:"guests"`
	RoomType         string    `json:"room_type" db:"room_type"`
	TotalPrice       float64   `json:"total_price" db:"total_price"`
	Currency         string    `json:"currency" db:"currency"`
	Status           string    `json:"status" db:"status"`
	BookingReference string    `json:"booking_reference" db:"booking_reference"`
	SpecialRequests  string    `json:"special_requests" db:"special_requests"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// HotelMapping mapea IDs internos con IDs de Amadeus
type HotelMapping struct {
	ID               int       `json:"id" db:"id"`
	InternalHotelID  string    `json:"internal_hotel_id" db:"internal_hotel_id"`
	AmadeusHotelID   string    `json:"amadeus_hotel_id" db:"amadeus_hotel_id"`
	HotelName        string    `json:"hotel_name" db:"hotel_name"`
	City             string    `json:"city" db:"city"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// AvailabilityCache representa datos de disponibilidad en caché
type AvailabilityCache struct {
	ID             int       `json:"id" db:"id"`
	HotelID        string    `json:"hotel_id" db:"hotel_id"`
	CheckInDate    time.Time `json:"check_in_date" db:"check_in_date"`
	CheckOutDate   time.Time `json:"check_out_date" db:"check_out_date"`
	Guests         int       `json:"guests" db:"guests"`
	Available      bool      `json:"available" db:"available"`
	Price          *float64  `json:"price" db:"price"`
	Currency       string    `json:"currency" db:"currency"`
	RoomsAvailable *int      `json:"rooms_available" db:"rooms_available"`
	CachedAt       time.Time `json:"cached_at" db:"cached_at"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
}

// Requests y Responses
type RegisterRequest struct {
	Email       string     `json:"email" validate:"required,email"`
	Password    string     `json:"password" validate:"required,min=6"`
	FirstName   string     `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string     `json:"last_name" validate:"required,min=2,max=50"`
	Role        UserRole   `json:"role"` // ← LÍNEA AGREGADA (opcional para registro)
	Phone       string     `json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AvailabilityRequest struct {
	HotelID      string    `json:"hotel_id" validate:"required"`
	CheckInDate  time.Time `json:"check_in_date" validate:"required"`
	CheckOutDate time.Time `json:"check_out_date" validate:"required"`
	Guests       int       `json:"guests" validate:"required,min=1,max=10"`
}

type CreateBookingRequest struct {
	HotelID         string    `json:"hotel_id" validate:"required"`
	CheckInDate     time.Time `json:"check_in_date" validate:"required"`
	CheckOutDate    time.Time `json:"check_out_date" validate:"required"`
	Guests          int       `json:"guests" validate:"required,min=1,max=10"`
	RoomType        string    `json:"room_type"`
	SpecialRequests string    `json:"special_requests"`
}

type AvailabilityResponse struct {
	HotelID        string   `json:"hotel_id"`
	Available      bool     `json:"available"`
	Price          *float64 `json:"price"`
	Currency       string   `json:"currency"`
	RoomsAvailable *int     `json:"rooms_available"`
	CheckInDate    string   `json:"check_in_date"`
	CheckOutDate   string   `json:"check_out_date"`
	Guests         int      `json:"guests"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// AmadeusTokenResponse para respuesta de autenticación de Amadeus
type AmadeusTokenResponse struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	Application string `json:"application"`
	ClientID    string `json:"client_id"`
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	State       string `json:"state"`
	Scope       string `json:"scope"`
}

// AmadeusHotelOffer representa una oferta de hotel de Amadeus
type AmadeusHotelOffer struct {
	Type         string  `json:"type"`
	Hotel        AmadeusHotel `json:"hotel"`
	Available    bool    `json:"available"`
	Offers       []AmadeusOffer `json:"offers"`
	Self         string  `json:"self"`
}

type AmadeusHotel struct {
	Type        string `json:"type"`
	HotelID     string `json:"hotelId"`
	ChainCode   string `json:"chainCode"`
	DupeID      string `json:"dupeId"`
	Name        string `json:"name"`
	Rating      string `json:"rating"`
	CityCode    string `json:"cityCode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type AmadeusOffer struct {
	ID            string        `json:"id"`
	CheckInDate   string        `json:"checkInDate"`
	CheckOutDate  string        `json:"checkOutDate"`
	RoomQuantity  int           `json:"roomQuantity"`
	RateCode      string        `json:"rateCode"`
	Room          AmadeusRoom   `json:"room"`
	Guests        AmadeusGuests `json:"guests"`
	Price         AmadeusPrice  `json:"price"`
}

type AmadeusRoom struct {
	Type        string `json:"type"`
	TypeCode    string `json:"typeCode"`
	Description string `json:"description"`
}

type AmadeusGuests struct {
	Adults int `json:"adults"`
}

type AmadeusPrice struct {
	Currency string `json:"currency"`
	Base     string `json:"base"`
	Total    string `json:"total"`
}