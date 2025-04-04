package persistence

import (
	"database/sql"
	"fmt"
	"strings" // Add this import

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



func (r *MySQLRepository) GetUserAlerts(userID int) (map[string]interface{}, error) {
	
	query := `SELECT numero_serie FROM ESP32 WHERE idUser = ?`
	
	fmt.Printf("Executing query for user ID %d: %s\n", userID, query)
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user ESP32 devices: %w", err)
	}
	defer rows.Close()
	
	// Collect all device serial numbers
	var serialNumbers []string
	for rows.Next() {
		var serialNumber string
		if err := rows.Scan(&serialNumber); err != nil {
			return nil, fmt.Errorf("error scanning ESP32 serial number: %w", err)
		}
		serialNumbers = append(serialNumbers, serialNumber)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ESP32 rows: %w", err)
	}
	
	fmt.Printf("Found %d devices for user ID %d: %v\n", len(serialNumbers), userID, serialNumbers)
	
	// If no devices found, return empty result
	if len(serialNumbers) == 0 {
		return map[string]interface{}{
			"message": "No devices found for this user",
			"alerts":  map[string]interface{}{},
		}, nil
	}
	
	// Now fetch alerts from all sensor tables for these serial numbers
	result := map[string]interface{}{
		"user_id": userID,
		"devices": serialNumbers,
		"alerts": map[string]interface{}{},
	}
	
	alertsMap := result["alerts"].(map[string]interface{})
	
	// Fetch KY_026 alerts
	ky026Alerts, err := r.fetchAlertsFromTable("KY_026", serialNumbers)
	if err != nil {
		return nil, err
	}
	alertsMap["KY_026"] = ky026Alerts
	
	// Fetch MQ_2 alerts
	mq2Alerts, err := r.fetchAlertsFromTable("MQ_2", serialNumbers)
	if err != nil {
		return nil, err
	}
	alertsMap["MQ_2"] = mq2Alerts
	
	// Fetch MQ_135 alerts
	mq135Alerts, err := r.fetchAlertsFromTable("MQ_135", serialNumbers)
	if err != nil {
		return nil, err
	}
	alertsMap["MQ_135"] = mq135Alerts
	
	// Fetch DHT_22 alerts
	dht22Alerts, err := r.fetchAlertsFromTable("DHT_22", serialNumbers)
	if err != nil {
		return nil, err
	}
	alertsMap["DHT_22"] = dht22Alerts
	
	return result, nil
}

// Helper method to fetch alerts from a specific table
func (r *MySQLRepository) fetchAlertsFromTable(tableName string, serialNumbers []string) ([]map[string]interface{}, error) {
	// Create placeholders for the IN clause
	placeholders := make([]string, len(serialNumbers))
	args := make([]interface{}, len(serialNumbers))
	for i, sn := range serialNumbers {
		placeholders[i] = "?"
		args[i] = sn
	}
	
	// Use the correct column names based on the table
	idColumn := fmt.Sprintf("id%s", tableName)
	
	query := fmt.Sprintf(
		`SELECT %s, fecha_activacion, fecha_desactivacion, estado, numero_serie 
		FROM %s 
		WHERE numero_serie IN (%s)
		ORDER BY fecha_activacion DESC`,
		idColumn,
		tableName,
		strings.Join(placeholders, ","),
	)
	
	fmt.Printf("Executing query for %s: %s\n", tableName, query)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching alerts from %s: %w", tableName, err)
	}
	defer rows.Close()
	
	var alerts []map[string]interface{}
	for rows.Next() {
		var id int
		var fechaActivacion, fechaDesactivacion, numeroSerie string
		var estado interface{}
		
		if err := rows.Scan(&id, &fechaActivacion, &fechaDesactivacion, &estado, &numeroSerie); err != nil {
			return nil, fmt.Errorf("error scanning alert from %s: %w", tableName, err)
		}
		
		alert := map[string]interface{}{
			"id":                  id,
			"fecha_activacion":    fechaActivacion,
			"fecha_desactivacion": fechaDesactivacion,
			"estado":              estado,
			"numero_serie":        numeroSerie,
		}
		
		alerts = append(alerts, alert)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating alerts from %s: %w", tableName, err)
	}
	
	return alerts, nil
}

// DB returns the database connection
func (r *MySQLRepository) DB() *sql.DB {
    return r.db
}