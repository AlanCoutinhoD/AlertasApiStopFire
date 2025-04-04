package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"hex_go/internal/domain/entities"
	"hex_go/pkg/config"
)

// RabbitMQClient handles the connection to RabbitMQ
type RabbitMQClient struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queueKY026   string
	queueMQ2     string
	queueMQ135   string
	queueDHT22   string
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient(cfg *config.Config) (*RabbitMQClient, error) {
	// Create connection string
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Create exchange
	err = channel.ExchangeDeclare(
		cfg.RabbitMQExchange, // name
		"direct",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Create queues
	queues := map[string]string{
		"KY026":  cfg.RabbitMQQueueKY026,
		"MQ2":    cfg.RabbitMQQueueMQ2,
		"MQ135":  cfg.RabbitMQQueueMQ135,
		"DHT22":  cfg.RabbitMQQueueDHT22,
	}

	for sensorType, queueName := range queues {
		_, err = channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		// Bind queue to exchange
		err = channel.QueueBind(
			queueName,            // queue name
			sensorType,           // routing key
			cfg.RabbitMQExchange, // exchange
			false,
			nil,
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	return &RabbitMQClient{
		conn:         conn,
		channel:      channel,
		exchangeName: cfg.RabbitMQExchange,
		queueKY026:   cfg.RabbitMQQueueKY026,
		queueMQ2:     cfg.RabbitMQQueueMQ2,
		queueMQ135:   cfg.RabbitMQQueueMQ135,
		queueDHT22:   cfg.RabbitMQQueueDHT22,
	}, nil
}

// PublishSensorData publishes sensor data to the appropriate queue
func (c *RabbitMQClient) PublishSensorData(data *entities.SensorDataRequest) error {
	// Convert data to JSON
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal sensor data: %w", err)
	}

	// Publish message
	err = c.channel.Publish(
		c.exchangeName, // exchange
		data.Sensor,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to %s queue: %s", data.Sensor, string(body))
	return nil
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}