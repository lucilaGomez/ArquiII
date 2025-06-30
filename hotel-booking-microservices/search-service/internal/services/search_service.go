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
	solrClient      *solr.Client
	rabbitConsumer  *rabbitmq.Consumer
	httpClient      *http.Client
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

	switch event.Type {
	case "created", "updated":
		// Indexar o actualizar hotel en Solr
		doc := solr.SolrDocument{
			ID:          event.Hotel.ID,
			Name:        event.Hotel.Name,
			Description: event.Hotel.Description,
			City:        event.Hotel.City,
			Address:     event.Hotel.Address,
			Amenities:   event.Hotel.Amenities,
			Rating:      event.Hotel.Rating,
			MinPrice:    event.Hotel.PriceRange.MinPrice,
			MaxPrice:    event.Hotel.PriceRange.MaxPrice,
			Currency:    event.Hotel.PriceRange.Currency,
			Thumbnail:   event.Hotel.Thumbnail,
			IsActive:    event.Hotel.IsActive,
		}

		err := s.solrClient.IndexHotel(doc)
		if err != nil {
			return fmt.Errorf("error indexing hotel %s in Solr: %v", event.HotelID, err)
		}

		log.Printf("✅ Hotel %s sincronizado en Solr", event.HotelID)

	case "deleted":
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
				"name":        hotel.Name,
				"description": hotel.Description,
				"city":        hotel.City,
				"rating":      hotel.Rating,
				"min_price":   hotel.MinPrice,
				"max_price":   hotel.MaxPrice,
				"currency":    hotel.Currency,
				"thumbnail":   hotel.Thumbnail,
				"amenities":   hotel.Amenities,
				"available":   true, // Asumimos disponible si no verificamos
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
				"name":        hotel.Name,
				"description": hotel.Description,
				"city":        hotel.City,
				"rating":      hotel.Rating,
				"min_price":   hotel.MinPrice,
				"max_price":   hotel.MaxPrice,
				"currency":    availability.Currency,
				"thumbnail":   hotel.Thumbnail,
				"amenities":   hotel.Amenities,
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
	// Llamar al booking-service para verificar disponibilidad
	url := fmt.Sprintf("http://booking-service:8080/api/v1/availability/%s?checkin=%s&checkout=%s&guests=%s",
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