package amadeus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"booking-service/internal/models"
)

// Client cliente para API de Amadeus
type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	accessToken  string
	tokenExpiry  time.Time
}

// NewClient crea un nuevo cliente de Amadeus
func NewClient(baseURL, clientID, clientSecret string) *Client {
	return &Client{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAccessToken obtiene un token de acceso de Amadeus
func (c *Client) GetAccessToken() error {
	// Si el token es válido y no ha expirado, no hacer nada
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-5*time.Minute)) {
		return nil
	}

	log.Println("🔑 Obteniendo nuevo token de Amadeus...")

	// Preparar datos del formulario
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	// Crear petición
	req, err := http.NewRequest("POST", c.baseURL+"/v1/security/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creando petición: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Ejecutar petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error ejecutando petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error obteniendo token: %d - %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta
	var tokenResp models.AmadeusTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return fmt.Errorf("error parseando respuesta: %v", err)
	}

	// Guardar token
	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	log.Printf("✅ Token de Amadeus obtenido, expira en %d segundos", tokenResp.ExpiresIn)
	return nil
}

// SearchHotelsByCity busca hoteles en una ciudad
func (c *Client) SearchHotelsByCity(cityCode string) ([]models.AmadeusHotel, error) {
	// Asegurar que tenemos token válido
	err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Construir URL
	endpoint := fmt.Sprintf("%s/v1/reference-data/locations/hotels/by-city?cityCode=%s", c.baseURL, cityCode)

	// Crear petición
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	// Ejecutar petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error buscando hoteles: %d - %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta
	var response struct {
		Data []models.AmadeusHotel `json:"data"`
	}
	
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	return response.Data, nil
}

// GetHotelOffers obtiene ofertas de un hotel específico
func (c *Client) GetHotelOffers(hotelID, checkInDate, checkOutDate string, adults int) ([]models.AmadeusHotelOffer, error) {
	// Asegurar que tenemos token válido
	err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Construir URL con parámetros
	params := url.Values{}
	params.Add("hotelIds", hotelID)
	params.Add("checkInDate", checkInDate)
	params.Add("checkOutDate", checkOutDate)
	params.Add("adults", fmt.Sprintf("%d", adults))

	endpoint := fmt.Sprintf("%s/v3/shopping/hotel-offers?%s", c.baseURL, params.Encode())

	// Crear petición
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	// Ejecutar petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error obteniendo ofertas: %d - %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta
	var response struct {
		Data []models.AmadeusHotelOffer `json:"data"`
	}
	
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	return response.Data, nil
}

// CreateBooking crea una reserva en Amadeus
func (c *Client) CreateBooking(offerID string, guestInfo map[string]interface{}) (string, error) {
	// Asegurar que tenemos token válido
	err := c.GetAccessToken()
	if err != nil {
		return "", err
	}

	// Preparar datos de la reserva
	bookingData := map[string]interface{}{
		"data": map[string]interface{}{
			"type":    "hotel-booking",
			"hotelId": offerID,
			"guests":  guestInfo,
		},
	}

	// Serializar datos
	jsonData, err := json.Marshal(bookingData)
	if err != nil {
		return "", fmt.Errorf("error serializando datos: %v", err)
	}

	// Crear petición
	endpoint := c.baseURL + "/v1/booking/hotel-bookings"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creando petición: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error ejecutando petición: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("error creando reserva: %d - %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta para obtener ID de reserva
	var response struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parseando respuesta: %v", err)
	}

	log.Printf("✅ Reserva creada en Amadeus: %s", response.Data.ID)
	return response.Data.ID, nil
}

// IsTokenValid verifica si el token actual es válido
func (c *Client) IsTokenValid() bool {
	return c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-5*time.Minute))
}