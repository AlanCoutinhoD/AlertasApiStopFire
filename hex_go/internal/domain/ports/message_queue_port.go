package ports

import "hex_go/internal/domain/entities"

type MessageQueuePort interface {
    PublishSensorData(data *entities.SensorDataRequest) error
    Close() error
}