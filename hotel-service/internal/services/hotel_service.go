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
	collection *mongo.Collection
	rabbit     *rabbitmq.Connection
}

// NewHotelService crea una nueva instancia del servicio
func NewHotelService(mongoClient *mongo.Client, rabbitConn *rabbitmq.Connection) *HotelService {
	collection := mongodb.GetCollection(mongoClient, "hotels_db", "hotels")
	
	return &HotelService{
		collection: collection,
		rabbit:     rabbitConn,
	}
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

// CreateHotel crea un nuevo hotel
func (s *HotelService) CreateHotel(req *models.CreateHotelRequest) (*models.Hotel, error) {
	// Convertir request a modelo
	hotel := req.ToHotel()

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insertar en MongoDB
	result, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, fmt.Errorf("error creando hotel: %v", err)
	}

	// Asignar el ID generado
	hotel.ID = result.InsertedID.(primitive.ObjectID)

	// Publicar evento en RabbitMQ
	event := models.HotelEvent{
		Type:      "created",
		HotelID:   hotel.ID,
		Hotel:     hotel,
		Timestamp: time.Now(),
	}

	err = s.rabbit.PublishEvent("hotel.created", event)
	if err != nil {
		// Log el error pero no fallar la operación
		fmt.Printf("Error publicando evento: %v\n", err)
	}

	return hotel, nil
}

// UpdateHotel actualiza un hotel existente
func (s *HotelService) UpdateHotel(id string, req *models.UpdateHotelRequest) (*models.Hotel, error) {
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

	// Agregar campos que no sean nil
	setFields := updateDoc["$set"].(bson.M)
	
	if req.Name != nil {
		setFields["name"] = *req.Name
	}
	if req.Description != nil {
		setFields["description"] = *req.Description
	}
	if req.City != nil {
		setFields["city"] = *req.City
	}
	if req.Address != nil {
		setFields["address"] = *req.Address
	}
	if req.Photos != nil {
		setFields["photos"] = req.Photos
	}
	if req.Thumbnail != nil {
		setFields["thumbnail"] = *req.Thumbnail
	}
	if req.Amenities != nil {
		setFields["amenities"] = req.Amenities
	}
	if req.Rating != nil {
		setFields["rating"] = *req.Rating
	}
	if req.PriceRange != nil {
		setFields["price_range"] = *req.PriceRange
	}
	if req.Contact != nil {
		setFields["contact"] = *req.Contact
	}
	if req.IsActive != nil {
		setFields["is_active"] = *req.IsActive
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
	event := models.HotelEvent{
		Type:      "updated",
		HotelID:   updatedHotel.ID,
		Hotel:     updatedHotel,
		Timestamp: time.Now(),
	}

	err = s.rabbit.PublishEvent("hotel.updated", event)
	if err != nil {
		// Log el error pero no fallar la operación
		fmt.Printf("Error publicando evento: %v\n", err)
	}

	return updatedHotel, nil
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