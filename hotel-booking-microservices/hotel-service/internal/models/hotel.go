package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Hotel representa un hotel en el sistema
type Hotel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	City        string            `json:"city" bson:"city"`
	Address     string            `json:"address" bson:"address"`
	Photos      []string          `json:"photos" bson:"photos"`
	Thumbnail   string            `json:"thumbnail" bson:"thumbnail"` // URL de la imagen principal
	Amenities   []string          `json:"amenities" bson:"amenities"`
	Rating      float64           `json:"rating" bson:"rating"`
	PriceRange  PriceRange        `json:"price_range" bson:"price_range"`
	Contact     Contact           `json:"contact" bson:"contact"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
	IsActive    bool              `json:"is_active" bson:"is_active"`
}

// PriceRange representa el rango de precios
type PriceRange struct {
	MinPrice int    `json:"min_price" bson:"min_price"`
	MaxPrice int    `json:"max_price" bson:"max_price"`
	Currency string `json:"currency" bson:"currency"`
}

// Contact representa la información de contacto
type Contact struct {
	Phone   string `json:"phone" bson:"phone"`
	Email   string `json:"email" bson:"email"`
	Website string `json:"website" bson:"website"`
}

// CreateHotelRequest representa la solicitud para crear un hotel
type CreateHotelRequest struct {
	Name        string     `json:"name" binding:"required" validate:"min=2,max=100"`
	Description string     `json:"description" binding:"required" validate:"min=10,max=1000"`
	City        string     `json:"city" binding:"required" validate:"min=2,max=50"`
	Address     string     `json:"address" binding:"required" validate:"min=10,max=200"`
	Amenities   []string   `json:"amenities"`
	Rating      float64    `json:"rating" validate:"min=1,max=5"`
	PriceRange  PriceRange `json:"price_range" binding:"required"`
	Contact     Contact    `json:"contact" binding:"required"`
	Thumbnail   string     `json:"thumbnail"` // URL de la imagen subida
}

// UpdateHotelRequest representa la solicitud para actualizar un hotel
type UpdateHotelRequest struct {
	Name        string     `json:"name,omitempty" validate:"min=2,max=100"`
	Description string     `json:"description,omitempty" validate:"min=10,max=1000"`
	City        string     `json:"city,omitempty" validate:"min=2,max=50"`
	Address     string     `json:"address,omitempty" validate:"min=10,max=200"`
	Amenities   []string   `json:"amenities,omitempty"`
	Rating      float64    `json:"rating,omitempty" validate:"min=1,max=5"`
	PriceRange  *PriceRange `json:"price_range,omitempty"`
	Contact     *Contact   `json:"contact,omitempty"`
	Thumbnail   string     `json:"thumbnail,omitempty"` // URL de la imagen
}

// SearchHotelRequest representa los parámetros de búsqueda
type SearchHotelRequest struct {
	City      string  `json:"city,omitempty"`
	MinPrice  int     `json:"min_price,omitempty"`
	MaxPrice  int     `json:"max_price,omitempty"`
	MinRating float64 `json:"min_rating,omitempty"`
	Amenities []string `json:"amenities,omitempty"`
}

// ImageUploadResponse representa la respuesta de subida de imagen
type ImageUploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}