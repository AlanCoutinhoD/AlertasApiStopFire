package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"hex_go/internal/domain/entities"
)

// MySQLRepository implements the SensorRepository interface
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new MySQL repository
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

// CreateKY026 inserts a new KY_026 sensor record
func (r *MySQLRepository) CreateKY026(sensor *entities.SensorKY026) error {
	query := `INSERT INTO KY_026 (fecha_activacion, fecha_desactivacion, estado, numero_serie) 
              VALUES (?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, sensor.FechaActivacion, sensor.FechaDesactivacion, sensor.Estado, sensor.NumeroSerie)
	if err != nil {
		return fmt.Errorf("error creating KY_026 sensor: %w", err)
	}
	
	return nil
}

// CreateMQ2 inserts a new MQ_2 sensor record
func (r *MySQLRepository) CreateMQ2(sensor *entities.SensorMQ2) error {
	query := `INSERT INTO MQ_2 (fecha_activacion, fecha_desactivacion, estado, numero_serie) 
              VALUES (?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, sensor.FechaActivacion, sensor.FechaDesactivacion, sensor.Estado, sensor.NumeroSerie)
	if err != nil {
		return fmt.Errorf("error creating MQ_2 sensor: %w", err)
	}
	
	return nil
}

// CreateMQ135 inserts a new MQ_135 sensor record
func (r *MySQLRepository) CreateMQ135(sensor *entities.SensorMQ135) error {
	query := `INSERT INTO MQ_135 (fecha_activacion, fecha_desactivacion, estado, numero_serie) 
              VALUES (?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, sensor.FechaActivacion, sensor.FechaDesactivacion, sensor.Estado, sensor.NumeroSerie)
	if err != nil {
		return fmt.Errorf("error creating MQ_135 sensor: %w", err)
	}
	
	return nil
}

// CreateDHT22 inserts a new DHT_22 sensor record
func (r *MySQLRepository) CreateDHT22(sensor *entities.SensorDHT22) error {
	query := `INSERT INTO DHT_22 (fecha_activacion, fecha_desactivacion, estado, numero_serie) 
              VALUES (?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, sensor.FechaActivacion, sensor.FechaDesactivacion, sensor.Estado, sensor.NumeroSerie)
	if err != nil {
		return fmt.Errorf("error creating DHT_22 sensor: %w", err)
	}
	
	return nil
}