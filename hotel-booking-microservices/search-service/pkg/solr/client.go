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

type SolrDocument struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	City        string   `json:"city"`
	Address     string   `json:"address"`
	Amenities   []string `json:"amenities"`
	Rating      float64  `json:"rating"`
	MinPrice    float64  `json:"min_price"`
	MaxPrice    float64  `json:"max_price"`
	Currency    string   `json:"currency"`
	Thumbnail   string   `json:"thumbnail"`
	IsActive    bool     `json:"is_active"`
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
			"name":        doc.Name,
			"description": doc.Description,
			"city":        doc.City,
			"address":     doc.Address,
			"amenities":   doc.Amenities,
			"rating":      doc.Rating,
			"min_price":   doc.MinPrice,
			"max_price":   doc.MaxPrice,
			"currency":    doc.Currency,
			"thumbnail":   doc.Thumbnail,
			"is_active":   doc.IsActive,
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