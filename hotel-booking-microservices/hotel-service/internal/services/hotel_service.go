package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"hotel-service/internal/models"
	"hotel-service/pkg/mongodb"
	"hotel-service/pkg/rabbitmq"
)

// HotelService maneja la lógica de negocio de hoteles
type HotelService struct {
	collection  *mongo.Collection
	rabbit      *rabbitmq.Connection
	mongoClient *mongo.Client  // ✅ AGREGADO para healthcheck
}

// NewHotelService crea una nueva instancia del servicio
func NewHotelService(mongoClient *mongo.Client, rabbitConn *rabbitmq.Connection) *HotelService {
	collection := mongodb.GetCollection(mongoClient, "hotels_db", "hotels")
	
	return &HotelService{
		collection:  collection,
		rabbit:      rabbitConn,
		mongoClient: mongoClient,  // ✅ AGREGADO
	}
}

// ✅ NUEVO: HealthCheck verifica la conexión a MongoDB
func (s *HotelService) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.mongoClient.Ping(ctx, nil)
}

// ✅ NUEVO: IsRabbitMQConnected verifica si RabbitMQ está conectado
func (s *HotelService) IsRabbitMQConnected() bool {
	if s.rabbit == nil {
		return false
	}
	return s.rabbit.IsConnected()
}

// GetHotelByID obtiene un hotel por su ID
func (s *HotelService) GetHotelByID(id string) (*models.Hotel, error) {
	// Convertir string a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID de hotel inválido: %v", err)
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Buscar en MongoDB
	var hotel models.Hotel
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&hotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("hotel no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo hotel: %v", err)
	}

	return &hotel, nil
}

// CreateHotel crea un nuevo hotel - VERSIÓN COMPATIBLE CON HANDLER
func (s *HotelService) CreateHotel(req interface{}) (*models.Hotel, error) {
	// Crear el hotel base
	hotel := &models.Hotel{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	// Mapear desde la request (que viene como interface{} desde el handler)
	if reqMap, ok := req.(map[string]interface{}); ok {
		if name, exists := reqMap["name"].(string); exists {
			hotel.Name = name
		}
		if description, exists := reqMap["description"].(string); exists {
			hotel.Description = description
		}
		if city, exists := reqMap["city"].(string); exists {
			hotel.City = city
		}
		if address, exists := reqMap["address"].(string); exists {
			hotel.Address = address
		}
		if amenities, exists := reqMap["amenities"].([]interface{}); exists {
			for _, amenity := range amenities {
				if amenityStr, ok := amenity.(string); ok {
					hotel.Amenities = append(hotel.Amenities, amenityStr)
				}
			}
		}
		if images, exists := reqMap["images"].([]interface{}); exists {
			for _, image := range images {
				if imageStr, ok := image.(string); ok {
					hotel.Photos = append(hotel.Photos, imageStr)
				}
			}
		}
		if thumbnail, exists := reqMap["thumbnail"].(string); exists {
			hotel.Thumbnail = thumbnail
		}
		if rating, exists := reqMap["rating"].(float64); exists {
			hotel.Rating = rating
		}
		if priceRange, exists := reqMap["price_range"].(map[string]interface{}); exists {
			if minPrice, ok := priceRange["min_price"].(float64); ok {
				hotel.PriceRange.MinPrice = int(minPrice)
			}
			if maxPrice, ok := priceRange["max_price"].(float64); ok {
				hotel.PriceRange.MaxPrice = int(maxPrice)
			}
			if currency, ok := priceRange["currency"].(string); ok {
				hotel.PriceRange.Currency = currency
			}
		}
		if contact, exists := reqMap["contact"].(map[string]interface{}); exists {
			if phone, ok := contact["phone"].(string); ok {
				hotel.Contact.Phone = phone
			}
			if email, ok := contact["email"].(string); ok {
				hotel.Contact.Email = email
			}
			if website, ok := contact["website"].(string); ok {
				hotel.Contact.Website = website
			}
		}
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insertar en MongoDB
	_, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, fmt.Errorf("error creando hotel: %v", err)
	}

	// Publicar evento en RabbitMQ (simplificado para evitar errores)
	err = s.publishSimpleEvent("hotel.created", hotel.ID.Hex())
	if err != nil {
		// Log el error pero no fallar la operación
		fmt.Printf("Error publicando evento: %v\n", err)
	}

	return hotel, nil
}

// UpdateHotel actualiza un hotel existente - VERSIÓN COMPATIBLE CON HANDLER
func (s *HotelService) UpdateHotel(id string, req interface{}) (*models.Hotel, error) {
	// Convertir string a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID de hotel inválido: %v", err)
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Construir documento de actualización
	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Agregar campos desde la request
	setFields := updateDoc["$set"].(bson.M)
	
	if reqMap, ok := req.(map[string]interface{}); ok {
		if name, exists := reqMap["name"].(string); exists && name != "" {
			setFields["name"] = name
		}
		if description, exists := reqMap["description"].(string); exists && description != "" {
			setFields["description"] = description
		}
		if city, exists := reqMap["city"].(string); exists && city != "" {
			setFields["city"] = city
		}
		if address, exists := reqMap["address"].(string); exists && address != "" {
			setFields["address"] = address
		}
		if amenities, exists := reqMap["amenities"].([]interface{}); exists {
			var amenitiesStr []string
			for _, amenity := range amenities {
				if amenityStr, ok := amenity.(string); ok {
					amenitiesStr = append(amenitiesStr, amenityStr)
				}
			}
			setFields["amenities"] = amenitiesStr
		}
		if images, exists := reqMap["images"].([]interface{}); exists {
			var imagesStr []string
			for _, image := range images {
				if imageStr, ok := image.(string); ok {
					imagesStr = append(imagesStr, imageStr)
				}
			}
			setFields["photos"] = imagesStr
		}
		if thumbnail, exists := reqMap["thumbnail"].(string); exists && thumbnail != "" {
			setFields["thumbnail"] = thumbnail
		}
		if rating, exists := reqMap["rating"].(float64); exists {
			setFields["rating"] = rating
		}
		if priceRange, exists := reqMap["price_range"].(map[string]interface{}); exists {
			priceRangeDoc := bson.M{}
			if minPrice, ok := priceRange["min_price"].(float64); ok {
				priceRangeDoc["min_price"] = int(minPrice)
			}
			if maxPrice, ok := priceRange["max_price"].(float64); ok {
				priceRangeDoc["max_price"] = int(maxPrice)
			}
			if currency, ok := priceRange["currency"].(string); ok {
				priceRangeDoc["currency"] = currency
			}
			setFields["price_range"] = priceRangeDoc
		}
		if contact, exists := reqMap["contact"].(map[string]interface{}); exists {
			contactDoc := bson.M{}
			if phone, ok := contact["phone"].(string); ok {
				contactDoc["phone"] = phone
			}
			if email, ok := contact["email"].(string); ok {
				contactDoc["email"] = email
			}
			if website, ok := contact["website"].(string); ok {
				contactDoc["website"] = website
			}
			setFields["contact"] = contactDoc
		}
	}

	// Actualizar en MongoDB
	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		updateDoc,
	)
	if err != nil {
		return nil, fmt.Errorf("error actualizando hotel: %v", err)
	}

	// Obtener hotel actualizado
	updatedHotel, err := s.GetHotelByID(id)
	if err != nil {
		return nil, err
	}

	// Publicar evento en RabbitMQ
	err = s.publishSimpleEvent("hotel.updated", id)
	if err != nil {
		// Log el error pero no fallar la operación
		fmt.Printf("Error publicando evento: %v\n", err)
	}

	return updatedHotel, nil
}

// DeleteHotel elimina un hotel (soft delete) - NUEVO MÉTODO
func (s *HotelService) DeleteHotel(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID de hotel inválido: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Soft delete - marcar como inactivo
	updateDoc := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		return fmt.Errorf("error eliminando hotel: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("hotel no encontrado")
	}

	// Publicar evento
	err = s.publishSimpleEvent("hotel.deleted", id)
	if err != nil {
		fmt.Printf("Error publicando evento: %v\n", err)
	}

	return nil
}

// GetAllHotels obtiene todos los hoteles (para testing)
func (s *HotelService) GetAllHotels() ([]*models.Hotel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar todos los hoteles activos
	cursor, err := s.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo hoteles: %v", err)
	}
	defer cursor.Close(ctx)

	var hotels []*models.Hotel
	for cursor.Next(ctx) {
		var hotel models.Hotel
		if err := cursor.Decode(&hotel); err != nil {
			return nil, fmt.Errorf("error decodificando hotel: %v", err)
		}
		hotels = append(hotels, &hotel)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error en cursor: %v", err)
	}

	return hotels, nil
}

// publishSimpleEvent publica un evento simple en RabbitMQ
func (s *HotelService) publishSimpleEvent(eventType, hotelID string) error {
	if s.rabbit == nil {
		return fmt.Errorf("conexión a RabbitMQ no disponible")
	}

	// Crear un evento simple como map
	event := map[string]interface{}{
		"type":      eventType,
		"hotel_id":  hotelID,
		"timestamp": time.Now(),
	}

	// Usar el método de publicación de tu paquete rabbitmq
	return s.rabbit.PublishEvent(eventType, event)
}