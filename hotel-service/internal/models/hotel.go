package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Hotel representa la estructura de un hotel en MongoDB
type Hotel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string            `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Description string            `json:"description" bson:"description" validate:"required,min=10,max=1000"`
	City        string            `json:"city" bson:"city" validate:"required,min=2,max=50"`
	Address     string            `json:"address" bson:"address" validate:"required,min=5,max=200"`
	Photos      []string          `json:"photos" bson:"photos"`
	Thumbnail   string            `json:"thumbnail" bson:"thumbnail"`
	Amenities   []string          `json:"amenities" bson:"amenities"`
	Rating      float64           `json:"rating" bson:"rating" validate:"gte=0,lte=5"`
	PriceRange  PriceRange        `json:"price_range" bson:"price_range"`
	Contact     Contact           `json:"contact" bson:"contact"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
	IsActive    bool              `json:"is_active" bson:"is_active"`
}

// PriceRange representa el rango de precios del hotel
type PriceRange struct {
	MinPrice float64 `json:"min_price" bson:"min_price" validate:"gte=0"`
	MaxPrice float64 `json:"max_price" bson:"max_price" validate:"gte=0"`
	Currency string  `json:"currency" bson:"currency" validate:"required,len=3"`
}

// Contact representa información de contacto del hotel
type Contact struct {
	Phone   string `json:"phone" bson:"phone"`
	Email   string `json:"email" bson:"email" validate:"email"`
	Website string `json:"website" bson:"website"`
}

// CreateHotelRequest estructura para crear hoteles
type CreateHotelRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=100"`
	Description string     `json:"description" validate:"required,min=10,max=1000"`
	City        string     `json:"city" validate:"required,min=2,max=50"`
	Address     string     `json:"address" validate:"required,min=5,max=200"`
	Photos      []string   `json:"photos"`
	Thumbnail   string     `json:"thumbnail"`
	Amenities   []string   `json:"amenities"`
	Rating      float64    `json:"rating" validate:"gte=0,lte=5"`
	PriceRange  PriceRange `json:"price_range"`
	Contact     Contact    `json:"contact"`
}

// UpdateHotelRequest estructura para actualizar hoteles
type UpdateHotelRequest struct {
	Name        *string     `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string     `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	City        *string     `json:"city,omitempty" validate:"omitempty,min=2,max=50"`
	Address     *string     `json:"address,omitempty" validate:"omitempty,min=5,max=200"`
	Photos      []string    `json:"photos,omitempty"`
	Thumbnail   *string     `json:"thumbnail,omitempty"`
	Amenities   []string    `json:"amenities,omitempty"`
	Rating      *float64    `json:"rating,omitempty" validate:"omitempty,gte=0,lte=5"`
	PriceRange  *PriceRange `json:"price_range,omitempty"`
	Contact     *Contact    `json:"contact,omitempty"`
	IsActive    *bool       `json:"is_active,omitempty"`
}

// ToHotel convierte CreateHotelRequest a Hotel
func (chr *CreateHotelRequest) ToHotel() *Hotel {
	now := time.Now()
	return &Hotel{
		Name:        chr.Name,
		Description: chr.Description,
		City:        chr.City,
		Address:     chr.Address,
		Photos:      chr.Photos,
		Thumbnail:   chr.Thumbnail,
		Amenities:   chr.Amenities,
		Rating:      chr.Rating,
		PriceRange:  chr.PriceRange,
		Contact:     chr.Contact,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
	}
}

// HotelEvent representa eventos para RabbitMQ
type HotelEvent struct {
	Type      string             `json:"type"`      // "created", "updated"
	HotelID   primitive.ObjectID `json:"hotel_id"`
	Hotel     *Hotel            `json:"hotel,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}