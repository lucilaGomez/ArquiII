package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type HotelEvent struct {
	Type      string    `json:"type"`
	HotelID   string    `json:"hotel_id"`
	Timestamp time.Time `json:"timestamp"`
	Hotel     struct {
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
	} `json:"hotel"`
}

func NewConsumer(uri string) (*Consumer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error opening channel: %v", err)
	}

	consumer := &Consumer{
		conn:    conn,
		channel: ch,
	}

	// Configurar exchange y cola
	if err := consumer.setup(); err != nil {
		consumer.Close()
		return nil, err
	}

	log.Println("✅ RabbitMQ Consumer conectado exitosamente")
	return consumer, nil
}

func (c *Consumer) setup() error {
	// Declarar exchange (debe coincidir con hotel-service)
	err := c.channel.ExchangeDeclare(
		"hotel_events", // nombre
		"topic",        // tipo
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring exchange: %v", err)
	}

	// Declarar cola específica para search-service
	_, err = c.channel.QueueDeclare(
		"search_service_queue", // nombre
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring queue: %v", err)
	}

	// Vincular cola al exchange para todos los eventos de hotel
	err = c.channel.QueueBind(
		"search_service_queue", // queue name
		"hotel.*",              // routing key pattern
		"hotel_events",         // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error binding queue: %v", err)
	}

	return nil
}

func (c *Consumer) StartConsuming(callback func(HotelEvent) error) error {
	// Configurar QoS para controlar la carga
	err := c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("error setting QoS: %v", err)
	}

	msgs, err := c.channel.Consume(
		"search_service_queue", // queue
		"search-service",       // consumer tag
		false,                  // auto-ack (false para manual ack)
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return fmt.Errorf("error consuming messages: %v", err)
	}

	log.Println("🐰 RabbitMQ Consumer iniciado, esperando mensajes...")

	go func() {
		for msg := range msgs {
			var event HotelEvent
			
			// Deserializar mensaje
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Error unmarshalling message: %v", err)
				msg.Nack(false, false) // No requeue message
				continue
			}

			log.Printf("📨 Evento recibido: %s para hotel %s", event.Type, event.HotelID)

			// Procesar evento
			if err := callback(event); err != nil {
				log.Printf("❌ Error procesando evento: %v", err)
				// En caso de error, reintentamos más tarde
				msg.Nack(false, true) // Requeue message
			} else {
				log.Printf("✅ Evento procesado exitosamente: %s", event.HotelID)
				msg.Ack(false) // Acknowledge message
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	log.Println("🔌 Cerrando conexión RabbitMQ Consumer...")
	
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Health check para verificar conectividad
func (c *Consumer) IsConnected() bool {
	return c.conn != nil && !c.conn.IsClosed()
}