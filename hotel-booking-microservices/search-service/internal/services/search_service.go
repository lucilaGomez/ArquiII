package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"search-service/pkg/rabbitmq"
	"search-service/pkg/solr"
)

type SearchService struct {
	solrClient     *solr.Client
	rabbitConsumer *rabbitmq.Consumer
	httpClient     *http.Client
}

type AvailabilityResult struct {
	HotelID   string  `json:"hotel_id"`
	Available bool    `json:"available"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
}

func NewSearchService(solrClient *solr.Client, rabbitConsumer *rabbitmq.Consumer) *SearchService {
	return &SearchService{
		solrClient:     solrClient,
		rabbitConsumer: rabbitConsumer,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *SearchService) StartEventConsumer() error {
	log.Println("🚀 Iniciando consumer de eventos de hoteles...")

	return s.rabbitConsumer.StartConsuming(func(event rabbitmq.HotelEvent) error {
		return s.syncHotelToSolr(event)
	})
}

func (s *SearchService) syncHotelToSolr(event rabbitmq.HotelEvent) error {
	log.Printf("🔄 Sincronizando hotel %s en Solr (evento: %s)", event.HotelID, event.Type)

	// CORREGIDO: Los eventos llegan con formato "hotel.created", "hotel.updated", etc.
	switch event.Type {
	case "hotel.created", "hotel.updated": // ✅ AGREGADO "hotel." prefix
		// PROBLEMA ADICIONAL: Si el evento no incluye datos completos del hotel,
		// necesitamos hacer GET al hotel-service según las consignas

		var hotelData struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			City        string   `json:"city"`
			Address     string   `json:"address"`
			Amenities   []string `json:"amenities"`
			Rating      float64  `json:"rating"`
			PriceRange  struct {
				MinPrice float64 `json:"min_price"`
				MaxPrice float64 `json:"max_price"`
				Currency string  `json:"currency"`
			} `json:"price_range"`
			Thumbnail string `json:"thumbnail"`
			IsActive  bool   `json:"is_active"`
		}

		// Si el evento no tiene datos completos del hotel, hacer GET al hotel-service
		if event.Hotel.ID == "" || event.Hotel.Name == "" {
			log.Printf("📞 Obteniendo datos completos del hotel %s desde hotel-service", event.HotelID)

			// Hacer GET al hotel-service según las consignas del proyecto
			hotelFromService, err := s.getHotelFromService(event.HotelID)
			if err != nil {
				return fmt.Errorf("error obteniendo hotel desde hotel-service: %v", err)
			}
			hotelData = hotelFromService
		} else {
			// Usar datos del evento
			hotelData.ID = event.Hotel.ID
			hotelData.Name = event.Hotel.Name
			hotelData.Description = event.Hotel.Description
			hotelData.City = event.Hotel.City
			hotelData.Address = event.Hotel.Address
			hotelData.Amenities = event.Hotel.Amenities
			hotelData.Rating = event.Hotel.Rating
			hotelData.PriceRange = event.Hotel.PriceRange
			hotelData.Thumbnail = event.Hotel.Thumbnail
			hotelData.IsActive = event.Hotel.IsActive
		}

		// Indexar o actualizar hotel en Solr
		doc := solr.SolrDocument{
			ID:          hotelData.ID,
			Name:        []string{hotelData.Name},                 // ✅ Convertir a slice
			Description: []string{hotelData.Description},          // ✅ Convertir a slice
			City:        []string{hotelData.City},                 // ✅ Convertir a slice
			Address:     []string{hotelData.Address},              // ✅ Convertir a slice
			Amenities:   hotelData.Amenities,                      // Ya es slice
			Rating:      []float64{hotelData.Rating},              // ✅ Convertir a slice
			MinPrice:    []float64{hotelData.PriceRange.MinPrice}, // ✅ Convertir a slice
			MaxPrice:    []float64{hotelData.PriceRange.MaxPrice}, // ✅ Convertir a slice
			Currency:    []string{hotelData.PriceRange.Currency},  // ✅ Convertir a slice
			Thumbnail:   []string{hotelData.Thumbnail},            // ✅ Convertir a slice
			IsActive:    []bool{hotelData.IsActive},               // ✅ Convertir a slice
		}

		err := s.solrClient.IndexHotel(doc)
		if err != nil {
			return fmt.Errorf("error indexing hotel %s in Solr: %v", event.HotelID, err)
		}

		log.Printf("✅ Hotel %s sincronizado en Solr", event.HotelID)

	case "hotel.deleted": // ✅ AGREGADO "hotel." prefix
		// Eliminar hotel de Solr
		err := s.solrClient.DeleteHotel(event.HotelID)
		if err != nil {
			return fmt.Errorf("error deleting hotel %s from Solr: %v", event.HotelID, err)
		}

		log.Printf("🗑️ Hotel %s eliminado de Solr", event.HotelID)

	default:
		log.Printf("⚠️ Tipo de evento desconocido: %s", event.Type)
	}

	return nil
}

// NUEVA FUNCIÓN: GET al hotel-service según las consignas
func (s *SearchService) getHotelFromService(hotelID string) (struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	City        string   `json:"city"`
	Address     string   `json:"address"`
	Amenities   []string `json:"amenities"`
	Rating      float64  `json:"rating"`
	PriceRange  struct {
		MinPrice float64 `json:"min_price"`
		MaxPrice float64 `json:"max_price"`
		Currency string  `json:"currency"`
	} `json:"price_range"`
	Thumbnail string `json:"thumbnail"`
	IsActive  bool   `json:"is_active"`
}, error) {
	var result struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		City        string   `json:"city"`
		Address     string   `json:"address"`
		Amenities   []string `json:"amenities"`
		Rating      float64  `json:"rating"`
		PriceRange  struct {
			MinPrice float64 `json:"min_price"`
			MaxPrice float64 `json:"max_price"`
			Currency string  `json:"currency"`
		} `json:"price_range"`
		Thumbnail string `json:"thumbnail"`
		IsActive  bool   `json:"is_active"`
	}

	// GET al hotel-service según consignas del proyecto
	url := fmt.Sprintf("http://hotel-service:8080/api/v1/hotels/%s", hotelID)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return result, fmt.Errorf("error calling hotel-service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("hotel-service returned status %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			City        string   `json:"city"`
			Address     string   `json:"address"`
			Amenities   []string `json:"amenities"`
			Rating      float64  `json:"rating"`
			PriceRange  struct {
				MinPrice float64 `json:"min_price"`
				MaxPrice float64 `json:"max_price"`
				Currency string  `json:"currency"`
			} `json:"price_range"`
			Thumbnail string `json:"thumbnail"`
			IsActive  bool   `json:"is_active"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return result, fmt.Errorf("error decoding hotel-service response: %v", err)
	}

	return response.Data, nil
}

func (s *SearchService) SearchHotelsWithAvailability(query, city, checkin, checkout, guests string) (map[string]interface{}, error) {
	log.Printf("🔍 Buscando hoteles: query='%s', city='%s', checkin='%s', checkout='%s', guests='%s'",
		query, city, checkin, checkout, guests)

	// 1. Buscar hoteles en Solr
	hotels, err := s.solrClient.SearchHotels(query, city)
	if err != nil {
		log.Printf("❌ Error buscando en Solr: %v", err)
		// Fallback: devolver respuesta vacía pero no fallar
		return map[string]interface{}{
			"hotels":      []map[string]interface{}{},
			"total":       0,
			"page":        1,
			"page_size":   20,
			"total_pages": 0,
			"message":     "Búsqueda en Solr falló, sin resultados",
		}, nil
	}

	log.Printf("📋 Encontrados %d hoteles en Solr", len(hotels))

	// 2. Si no hay fechas, devolver solo los hoteles sin verificar disponibilidad
	if checkin == "" || checkout == "" {
		var hotelResults []map[string]interface{}
		for _, hotel := range hotels {
			hotelData := map[string]interface{}{
				"id":          hotel.ID,
				"name":        hotel.GetName(),        // ✅ Usar método helper
				"description": hotel.GetDescription(), // ✅ Usar método helper
				"city":        hotel.GetCity(),        // ✅ Usar método helper
				"rating":      hotel.GetRating(),      // ✅ Usar método helper
				"min_price":   hotel.GetMinPrice(),    // ✅ Usar método helper
				"max_price":   hotel.GetMaxPrice(),    // ✅ Usar método helper
				"currency":    hotel.GetCurrency(),    // ✅ Usar método helper
				"thumbnail":   hotel.GetThumbnail(),   // ✅ Usar método helper
				"amenities":   hotel.GetAmenities(),   // ✅ Usar método helper
				"available":   true,                   // Asumimos disponible si no verificamos
			}
			hotelResults = append(hotelResults, hotelData)
		}

		return map[string]interface{}{
			"hotels":      hotelResults,
			"total":       len(hotelResults),
			"page":        1,
			"page_size":   20,
			"total_pages": 1,
		}, nil
	}

	// 3. Calcular disponibilidad de forma concurrente
	log.Println("⚡ Verificando disponibilidad de forma concurrente...")
	availabilities := s.checkAvailabilityConcurrent(hotels, checkin, checkout, guests)

	// 4. Combinar resultados (solo hoteles con disponibilidad)
	var filteredHotels []map[string]interface{}
	for _, hotel := range hotels {
		availability, exists := availabilities[hotel.ID]
		if !exists {
			// Si no obtuvimos disponibilidad, asumimos no disponible
			log.Printf("⚠️ No se pudo verificar disponibilidad para hotel %s", hotel.ID)
			continue
		}

		// Solo incluir hoteles disponibles
		if availability.Available {
			hotelData := map[string]interface{}{
				"id":          hotel.ID,
				"name":        hotel.GetName(),        // ✅ Usar método helper
				"description": hotel.GetDescription(), // ✅ Usar método helper
				"city":        hotel.GetCity(),        // ✅ Usar método helper
				"rating":      hotel.GetRating(),      // ✅ Usar método helper
				"min_price":   hotel.GetMinPrice(),    // ✅ Usar método helper
				"max_price":   hotel.GetMaxPrice(),    // ✅ Usar método helper
				"currency":    availability.Currency,
				"thumbnail":   hotel.GetThumbnail(), // ✅ Usar método helper
				"amenities":   hotel.GetAmenities(), // ✅ Usar método helper
				"available":   availability.Available,
				"price":       availability.Price,
			}
			filteredHotels = append(filteredHotels, hotelData)
		}
	}

	log.Printf("✅ Devolviendo %d hoteles con disponibilidad confirmada", len(filteredHotels))

	return map[string]interface{}{
		"hotels":      filteredHotels,
		"total":       len(filteredHotels),
		"page":        1,
		"page_size":   20,
		"total_pages": 1,
	}, nil
}

// checkAvailabilityConcurrent consulta disponibilidad de forma concurrente
func (s *SearchService) checkAvailabilityConcurrent(hotels []solr.SolrDocument, checkin, checkout, guests string) map[string]AvailabilityResult {
	results := make(map[string]AvailabilityResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Canal para limitar concurrencia (máximo 10 consultas simultáneas)
	semaphore := make(chan struct{}, 10)

	log.Printf("🔄 Iniciando verificación concurrente para %d hoteles", len(hotels))

	for _, hotel := range hotels {
		wg.Add(1)
		go func(h solr.SolrDocument) {
			defer wg.Done()
			semaphore <- struct{}{}        // Adquirir semáforo
			defer func() { <-semaphore }() // Liberar semáforo

			availability := s.checkSingleAvailability(h.ID, checkin, checkout, guests)

			mu.Lock()
			results[h.ID] = availability
			mu.Unlock()
		}(hotel)
	}

	wg.Wait()
	log.Printf("✅ Verificación concurrente completada")
	return results
}

func (s *SearchService) checkSingleAvailability(hotelID, checkin, checkout, guests string) AvailabilityResult {
	// ✅ CORREGIDO: Puerto interno 8080, path /api/availability y parámetros en minúscula
	url := fmt.Sprintf("http://booking-service:8080/api/availability/%s?checkin=%s&checkout=%s&guests=%s",
		hotelID, checkin, checkout, guests)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("❌ Error creando request para hotel %s: %v", hotelID, err)
		return AvailabilityResult{HotelID: hotelID, Available: false}
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("❌ Error consultando disponibilidad hotel %s: %v", hotelID, err)
		return AvailabilityResult{HotelID: hotelID, Available: false}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ Booking service respondió %d para hotel %s", resp.StatusCode, hotelID)
		return AvailabilityResult{HotelID: hotelID, Available: false}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Error leyendo respuesta para hotel %s: %v", hotelID, err)
		return AvailabilityResult{HotelID: hotelID, Available: false}
	}

	var result struct {
		Data struct {
			Available bool    `json:"available"`
			Price     float64 `json:"price"`
			Currency  string  `json:"currency"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("❌ Error parseando respuesta para hotel %s: %v", hotelID, err)
		return AvailabilityResult{HotelID: hotelID, Available: false}
	}

	availability := AvailabilityResult{
		HotelID:   hotelID,
		Available: result.Data.Available,
		Price:     result.Data.Price,
		Currency:  result.Data.Currency,
	}

	if availability.Available {
		log.Printf("✅ Hotel %s disponible - precio: %.2f %s", hotelID, availability.Price, availability.Currency)
	} else {
		log.Printf("❌ Hotel %s no disponible", hotelID)
	}

	return availability
}

// Health check para Solr
func (s *SearchService) HealthCheck() map[string]interface{} {
	health := map[string]interface{}{
		"solr":     "unknown",
		"rabbitmq": "unknown",
	}

	// Check Solr
	if err := s.solrClient.Ping(); err != nil {
		health["solr"] = "down"
		log.Printf("❌ Solr health check failed: %v", err)
	} else {
		health["solr"] = "up"
	}

	// Check RabbitMQ
	if s.rabbitConsumer.IsConnected() {
		health["rabbitmq"] = "up"
	} else {
		health["rabbitmq"] = "down"
	}

	return health
}
