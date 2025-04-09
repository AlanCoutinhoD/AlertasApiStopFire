package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"hex_go/internal/domain/entities"
	"hex_go/internal/domain/ports"
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
		"KY_026": cfg.RabbitMQQueueKY026,
		"MQ_2":   cfg.RabbitMQQueueMQ2,
		"MQ_135": cfg.RabbitMQQueueMQ135,
		"DHT_22": cfg.RabbitMQQueueDHT22,
	}

	log.Printf("Creating and binding queues to exchange: %s", cfg.RabbitMQExchange)
	for sensorType, queueName := range queues {
		log.Printf("Declaring queue: %s for sensor type: %s", queueName, sensorType)
		
		// Declare the queue with more durable settings
		_, err = channel.QueueDeclare(
			queueName, // name
			true,      // durable - survive broker restart
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
		
		log.Printf("Queue %s declared successfully", queueName)

		// Bind queue to exchange
		log.Printf("Binding queue %s to exchange %s with routing key %s", 
			queueName, cfg.RabbitMQExchange, sensorType)
		
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
		
		log.Printf("Queue %s bound successfully to exchange %s", queueName, cfg.RabbitMQExchange)
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

	// Use sensor type directly as routing key
	routingKey := data.Sensor

	// Publish message
	err = c.channel.Publish(
		c.exchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to %s queue with routing key %s: %s", 
		getQueueNameForSensor(routingKey, c), routingKey, string(body))
	return nil
}


func getQueueNameForSensor(sensorType string, c *RabbitMQClient) string {
	switch sensorType {
	case "KY_026":
		return c.queueKY026
	case "MQ_2":
		return c.queueMQ2
	case "MQ_135":
		return c.queueMQ135
	case "DHT_22":
		return c.queueDHT22
	default:
		return "unknown"
	}
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitMQClient) Close() error {
    var err error
    if c.channel != nil {
        if cerr := c.channel.Close(); cerr != nil {
            err = fmt.Errorf("error closing channel: %v", cerr)
        }
    }
    if c.conn != nil {
        if cerr := c.conn.Close(); cerr != nil {
            if err != nil {
                err = fmt.Errorf("%v; error closing connection: %v", err, cerr)
            } else {
                err = fmt.Errorf("error closing connection: %v", cerr)
            }
        }
    }
    return err
}

// Verify interface implementation
var _ ports.MessageQueuePort = (*RabbitMQClient)(nil)