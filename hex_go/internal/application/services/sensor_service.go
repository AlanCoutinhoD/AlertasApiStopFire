package services

import (
	"errors"

	"hex_go/internal/domain/entities"
	"hex_go/internal/domain/repositories"
)

// SensorService handles the business logic for sensors
type SensorService struct {
	repo repositories.SensorRepository
}

// NewSensorService creates a new sensor service
func NewSensorService(repo repositories.SensorRepository) *SensorService {
	return &SensorService{
		repo: repo,
	}
}

// ProcessSensorData processes incoming sensor data and stores it in the appropriate table
func (s *SensorService) ProcessSensorData(data *entities.SensorDataRequest) error {
	switch data.Sensor {
	case "KY_026":
		sensor := &entities.SensorKY026{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		return s.repo.CreateKY026(sensor)
	case "MQ_2":
		sensor := &entities.SensorMQ2{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		return s.repo.CreateMQ2(sensor)
	case "MQ_135":
		sensor := &entities.SensorMQ135{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		return s.repo.CreateMQ135(sensor)
	case "DHT_22":
		sensor := &entities.SensorDHT22{
			FechaActivacion:    data.FechaActivacion,
			FechaDesactivacion: data.FechaDesactivacion,
			Estado:             data.Estado,
			NumeroSerie:        data.NumeroSerie,
		}
		return s.repo.CreateDHT22(sensor)
	default:
		return errors.New("sensor type not supported")
	}
}