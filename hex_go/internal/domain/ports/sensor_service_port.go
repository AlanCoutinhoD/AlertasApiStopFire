package ports

import "hex_go/internal/domain/entities"

type SensorServicePort interface {
    ProcessSensorData(data *entities.SensorDataRequest) error
    GetUserAlerts(userID int) (map[string]interface{}, error)
}