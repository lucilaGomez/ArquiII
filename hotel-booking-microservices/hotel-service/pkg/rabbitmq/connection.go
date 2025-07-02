package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection wrapper para RabbitMQ
type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect establece conexión con RabbitMQ
func Connect(uri string) (*Connection, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	// Crear canal
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error creando canal: %v", err)
	}

	// Crear instancia de conexión
	rabbitConn := &Connection{
		conn:    conn,
		channel: ch,
	}

	// Declarar exchange y cola
	err = rabbitConn.setupExchangeAndQueue()
	if err != nil {
		rabbitConn.Close()
		return nil, err
	}

	log.Println("✅ Conectado a RabbitMQ exitosamente")
	return rabbitConn, nil
}

// setupExchangeAndQueue configura exchange y cola
func (r *Connection) setupExchangeAndQueue() error {
	// Declarar exchange
	err := r.channel.ExchangeDeclare(
		"hotel_events", // nombre
		"topic",        // tipo
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("error declarando exchange: %v", err)
	}

	// Declarar cola
	_, err = r.channel.QueueDeclare(
		"hotel_updates", // nombre
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("error declarando cola: %v", err)
	}

	// Vincular cola al exchange
	err = r.channel.QueueBind(
		"hotel_updates",   // queue name
		"hotel.*",         // routing key
		"hotel_events",    // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error vinculando cola: %v", err)
	}

	return nil
}

// PublishEvent publica un evento en RabbitMQ
func (r *Connection) PublishEvent(routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshalling event: %v", err)
	}

	err = r.channel.Publish(
		"hotel_events", // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("error publicando evento: %v", err)
	}

	log.Printf("📤 Evento publicado: %s", routingKey)
	return nil
}

// Close cierra la conexión
func (r *Connection) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// IsConnected verifica si la conexión está activa
func (r *Connection) IsConnected() bool {
	if r.conn == nil || r.conn.IsClosed() {
		return false
	}
	if r.channel == nil {
		return false
	}
	return true
}
