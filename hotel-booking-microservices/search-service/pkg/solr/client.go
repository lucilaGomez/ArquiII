package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ✅ STRUCT CORREGIDO - Todos los campos como arrays porque así los devuelve Solr
type SolrDocument struct {
	ID          string    `json:"id"`
	Name        []string  `json:"name"`        // ✅ Cambio de string a []string
	Description []string  `json:"description"` // ✅ Cambio de string a []string
	City        []string  `json:"city"`        // ✅ Cambio de string a []string
	Address     []string  `json:"address"`     // ✅ Cambio de string a []string
	Amenities   []string  `json:"amenities"`   // ✅ Ya era []string
	Rating      []float64 `json:"rating"`      // ✅ Cambio de float64 a []float64
	MinPrice    []float64 `json:"min_price"`   // ✅ Cambio de float64 a []float64
	MaxPrice    []float64 `json:"max_price"`   // ✅ Cambio de float64 a []float64
	Currency    []string  `json:"currency"`    // ✅ Cambio de string a []string
	Thumbnail   []string  `json:"thumbnail"`   // ✅ Cambio de string a []string
	IsActive    []bool    `json:"is_active"`   // ✅ Cambio de bool a []bool
}

// ✅ MÉTODOS HELPER para obtener el primer valor de cada array
func (s *SolrDocument) GetName() string {
	if len(s.Name) > 0 {
		return s.Name[0]
	}
	return ""
}

func (s *SolrDocument) GetDescription() string {
	if len(s.Description) > 0 {
		return s.Description[0]
	}
	return ""
}

func (s *SolrDocument) GetCity() string {
	if len(s.City) > 0 {
		return s.City[0]
	}
	return ""
}

func (s *SolrDocument) GetAddress() string {
	if len(s.Address) > 0 {
		return s.Address[0]
	}
	return ""
}

func (s *SolrDocument) GetRating() float64 {
	if len(s.Rating) > 0 {
		return s.Rating[0]
	}
	return 0
}

func (s *SolrDocument) GetMinPrice() float64 {
	if len(s.MinPrice) > 0 {
		return s.MinPrice[0]
	}
	return 0
}

func (s *SolrDocument) GetMaxPrice() float64 {
	if len(s.MaxPrice) > 0 {
		return s.MaxPrice[0]
	}
	return 0
}

func (s *SolrDocument) GetCurrency() string {
	if len(s.Currency) > 0 {
		return s.Currency[0]
	}
	return ""
}

func (s *SolrDocument) GetThumbnail() string {
	if len(s.Thumbnail) > 0 {
		return s.Thumbnail[0]
	}
	return ""
}

func (s *SolrDocument) GetIsActive() bool {
	if len(s.IsActive) > 0 {
		return s.IsActive[0]
	}
	return false
}

// ✅ Método para obtener amenities (ya era slice, pero agregamos validación)
func (s *SolrDocument) GetAmenities() []string {
	if len(s.Amenities) > 0 {
		return s.Amenities
	}
	return []string{}
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) IndexHotel(doc SolrDocument) error {
	// Preparar documento para Solr
	solrDoc := []map[string]interface{}{
		{
			"id":          doc.ID,
			"name":        doc.GetName(), // ✅ Usar métodos helper
			"description": doc.GetDescription(),
			"city":        doc.GetCity(),
			"address":     doc.GetAddress(),
			"amenities":   doc.GetAmenities(),
			"rating":      doc.GetRating(),
			"min_price":   doc.GetMinPrice(),
			"max_price":   doc.GetMaxPrice(),
			"currency":    doc.GetCurrency(),
			"thumbnail":   doc.GetThumbnail(),
			"is_active":   doc.GetIsActive(),
		},
	}

	jsonData, err := json.Marshal(solrDoc)
	if err != nil {
		return fmt.Errorf("error marshalling document: %v", err)
	}

	// Enviar a Solr
	endpoint := fmt.Sprintf("%s/update?commit=true", c.baseURL)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) SearchHotels(query, city string) ([]SolrDocument, error) {
	searchQuery := "*:*"

	if query != "" {
		searchQuery = fmt.Sprintf("name:*%s* OR description:*%s*", query, query)
	}

	if city != "" {
		cityQuery := fmt.Sprintf("city:*%s*", city)
		if searchQuery == "*:*" {
			searchQuery = cityQuery
		} else {
			searchQuery = fmt.Sprintf("(%s) AND (%s)", searchQuery, cityQuery)
		}
	}

	// URL encode the query
	encodedQuery := url.QueryEscape(searchQuery)
	endpoint := fmt.Sprintf("%s/select?q=%s&wt=json&rows=100&fq=is_active:true", c.baseURL, encodedQuery)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error searching solr: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("solr search error: %d - %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response struct {
			Docs []SolrDocument `json:"docs"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding solr response: %v", err)
	}

	return result.Response.Docs, nil
}

func (c *Client) DeleteHotel(hotelID string) error {
	deleteDoc := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": hotelID,
		},
	}

	jsonData, err := json.Marshal(deleteDoc)
	if err != nil {
		return fmt.Errorf("error marshalling delete request: %v", err)
	}

	endpoint := fmt.Sprintf("%s/update?commit=true", c.baseURL)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating delete request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("solr delete error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) Ping() error {
	endpoint := fmt.Sprintf("%s/admin/ping", c.baseURL)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("error pinging solr: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr ping failed: %d", resp.StatusCode)
	}

	return nil
}
