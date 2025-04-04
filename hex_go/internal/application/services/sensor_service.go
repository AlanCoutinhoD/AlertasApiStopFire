package services

import (
	"errors"
	"fmt"

	"hex_go/internal/domain/entities"
	"hex_go/internal/domain/repositories"
	"hex_go/pkg/rabbitmq"
)

// SensorService handles the business logic for sensors
type SensorService struct {
	repo        repositories.SensorRepository
	rabbitClient *rabbitmq.RabbitMQClient
}

// NewSensorService creates a new sensor service
func NewSensorService(repo repositories.SensorRepository, rabbitClient *rabbitmq.RabbitMQClient) *SensorService {
	return &SensorService{
		repo:        repo,
		rabbitClient: rabbitClient,
	}
}

// ProcessSensorData processes incoming sensor data, stores it in the database, and publishes to RabbitMQ
func (s *SensorService) ProcessSensorData(data *entities.SensorDataRequest) error {
	// First, store in database
	var err error
	
	switch data.Sensor {
	case "KY_026":
		sensor := &entities.SensorKY026{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		err = s.repo.CreateKY026(sensor)
	case "MQ_2":
		sensor := &entities.SensorMQ2{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		err = s.repo.CreateMQ2(sensor)
	case "MQ_135":
		sensor := &entities.SensorMQ135{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		err = s.repo.CreateMQ135(sensor)
	case "DHT_22":
		// For DHT_22, convert estado to string
		estadoStr := fmt.Sprintf("%v", data.Estado)
		
		sensor := &entities.SensorDHT22{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             estadoStr,
			NumeroSerie:        data.NumeroSerie,
		}
		err = s.repo.CreateDHT22(sensor)
	default:
		return errors.New("sensor type not supported")
	}
	
	if err != nil {
		return err
	}
	
	// Then, publish to RabbitMQ
	if s.rabbitClient != nil {
		return s.rabbitClient.PublishSensorData(data)
	}
	
	return nil
}