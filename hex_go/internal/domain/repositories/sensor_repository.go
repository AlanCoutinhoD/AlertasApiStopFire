package repositories

import (
	"database/sql"
	"hex_go/internal/domain/entities"
)

// SensorRepository defines the interface for sensor data operations
type SensorRepository interface {
	CreateKY026(sensor *entities.SensorKY026) error
	CreateMQ2(sensor *entities.SensorMQ2) error
	CreateMQ135(sensor *entities.SensorMQ135) error
	CreateDHT22(sensor *entities.SensorDHT22) error

	// GetUserAlerts retrieves all alerts for a user based on their ID
	GetUserAlerts(userID int) (map[string]interface{}, error)

	DB() *sql.DB
}