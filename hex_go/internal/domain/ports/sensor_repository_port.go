package ports

import (
    "database/sql"
    "hex_go/internal/domain/entities"
)

type SensorRepositoryPort interface {
    CreateKY026(sensor *entities.SensorKY026) error
    CreateMQ2(sensor *entities.SensorMQ2) error
    CreateMQ135(sensor *entities.SensorMQ135) error
    CreateDHT22(sensor *entities.SensorDHT22) error
    GetUserAlerts(userID int) (map[string]interface{}, error)
    DB() *sql.DB
}