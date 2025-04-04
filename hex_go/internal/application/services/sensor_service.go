package services

import (
	"errors"
	"fmt"
	"strconv"
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
	    // Convert estado to int for KY_026
	    var estado int
	    switch v := data.Estado.(type) {
	    case float64:
	        estado = int(v)
	    case int:
	        estado = v
	    case string:
	        var parseErr error
	        estado, parseErr = strconv.Atoi(v)
	        if parseErr != nil {
	            return fmt.Errorf("invalid estado value for KY_026: %v", parseErr)
	        }
	    default:
	        return fmt.Errorf("invalid estado type for KY_026: %T", v)
	    }
	    
	    sensor := &entities.SensorKY026{
	        FechaActivacion:    data.FechaActivacion,
	        FechaDesactivacion: data.FechaDesactivacion,
	        Estado:             estado,
	        NumeroSerie:        data.NumeroSerie,
	    }
	    err = s.repo.CreateKY026(sensor)
	case "MQ_2":
	    // Convert estado to int for MQ_2
	    var estado int
	    switch v := data.Estado.(type) {
	    case float64:
	        estado = int(v)
	    case int:
	        estado = v
	    case string:
	        var parseErr error
	        estado, parseErr = strconv.Atoi(v)
	        if parseErr != nil {
	            return fmt.Errorf("invalid estado value for MQ_2: %v", parseErr)
	        }
	    default:
	        return fmt.Errorf("invalid estado type for MQ_2: %T", v)
	    }
	    
	    sensor := &entities.SensorMQ2{
	        FechaActivacion:    data.FechaActivacion,
	        FechaDesactivacion: data.FechaDesactivacion,
	        Estado:             estado,
	        NumeroSerie:        data.NumeroSerie,
	    }
	    err = s.repo.CreateMQ2(sensor)
	case "MQ_135":
	    // Convert estado to int for MQ_135
	    var estado int
	    switch v := data.Estado.(type) {
	    case float64:
	        estado = int(v)
	    case int:
	        estado = v
	    case string:
	        var parseErr error
	        estado, parseErr = strconv.Atoi(v)
	        if parseErr != nil {
	            return fmt.Errorf("invalid estado value for MQ_135: %v", parseErr)
	        }
	    default:
	        return fmt.Errorf("invalid estado type for MQ_135: %T", v)
	    }
	    
	    sensor := &entities.SensorMQ135{
	        FechaActivacion:    data.FechaActivacion,
	        FechaDesactivacion: data.FechaDesactivacion,
	        Estado:             estado,
	        NumeroSerie:        data.NumeroSerie,
	    }
	    err = s.repo.CreateMQ135(sensor)
	case "DHT_22":
	    // For DHT_22, convert estado to string
	    var estadoStr string
	    switch v := data.Estado.(type) {
	    case string:
	        estadoStr = v
	    case float64:
	        estadoStr = fmt.Sprintf("%v", v)
	    case int:
	        estadoStr = fmt.Sprintf("%d", v)
	    default:
	        estadoStr = fmt.Sprintf("%v", v)
	    }
	    
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