package memcached

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Client wrapper para Memcached
type Client struct {
	mc *memcache.Client
}

// Connect establece conexión con Memcached
func Connect(uri string) (*Client, error) {
	// Crear cliente Memcached
	mc := memcache.New(uri)
	
	// Configurar timeouts
	mc.Timeout = 100 * time.Millisecond
	mc.MaxIdleConns = 100

	// Verificar conexión
	err := mc.Ping()
	if err != nil {
		return nil, fmt.Errorf("error conectando a Memcached: %v", err)
	}

	log.Println("✅ Conectado a Memcached exitosamente")
	
	return &Client{mc: mc}, nil
}

// Set almacena un valor en caché
func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	// Serializar valor
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error serializando valor: %v", err)
	}

	// Crear item
	item := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(expiration.Seconds()),
	}

	// Almacenar en Memcached
	err = c.mc.Set(item)
	if err != nil {
		return fmt.Errorf("error almacenando en caché: %v", err)
	}

	log.Printf("💾 Valor almacenado en caché: %s", key)
	return nil
}

// Get obtiene un valor del caché
func (c *Client) Get(key string, result interface{}) error {
	// Obtener item
	item, err := c.mc.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("error obteniendo de caché: %v", err)
	}

	// Deserializar valor
	err = json.Unmarshal(item.Value, result)
	if err != nil {
		return fmt.Errorf("error deserializando valor: %v", err)
	}

	log.Printf("🔍 Valor obtenido del caché: %s", key)
	return nil
}

// Delete elimina un valor del caché
func (c *Client) Delete(key string) error {
	err := c.mc.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("error eliminando de caché: %v", err)
	}

	log.Printf("🗑️ Valor eliminado del caché: %s", key)
	return nil
}

// Exists verifica si una clave existe en el caché
func (c *Client) Exists(key string) bool {
	_, err := c.mc.Get(key)
	return err == nil
}

// Close cierra la conexión (Memcached no requiere cierre explícito)
func (c *Client) Close() error {
	// Memcached client no requiere cierre explícito
	return nil
}

// GenerateAvailabilityKey genera una clave única para disponibilidad
func GenerateAvailabilityKey(hotelID string, checkIn, checkOut time.Time, guests int) string {
	return fmt.Sprintf("availability:%s:%s:%s:%d",
		hotelID,
		checkIn.Format("2006-01-02"),
		checkOut.Format("2006-01-02"),
		guests,
	)
}

// GenerateTokenKey genera una clave para tokens de Amadeus
func GenerateTokenKey() string {
	return "amadeus:token"
}