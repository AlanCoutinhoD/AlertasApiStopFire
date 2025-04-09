package entities

// SensorKY026 represents the KY_026 sensor entity
type SensorKY026 struct {
	ID                 int    `json:"id"`
	FechaActivacion    string `json:"fecha_activacion"`
	FechaDesactivacion string `json:"fecha_desactivacion"`
	Estado             int    `json:"estado"`
	NumeroSerie        string `json:"numero_serie"`
}

// SensorMQ2 represents the MQ_2 sensor entity
type SensorMQ2 struct {
	ID                 int    `json:"id"`
	FechaActivacion    string `json:"fecha_activacion"`
	FechaDesactivacion string `json:"fecha_desactivacion"`
	Estado             int    `json:"estado"`
	NumeroSerie        string `json:"numero_serie"`
}

// SensorMQ135 represents the MQ_135 sensor entity
type SensorMQ135 struct {
	ID                 int    `json:"id"`
	FechaActivacion    string `json:"fecha_activacion"`
	FechaDesactivacion string `json:"fecha_desactivacion"`
	Estado             int    `json:"estado"`
	NumeroSerie        string `json:"numero_serie"`
}

// SensorDHT22 represents the DHT_22 sensor entity
type SensorDHT22 struct {
	ID                 int    `json:"id"`
	FechaActivacion    string `json:"fecha_activacion"`
	FechaDesactivacion string `json:"fecha_desactivacion"`
	Estado             string `json:"estado"` // Changed from int to string
	NumeroSerie        string `json:"numero_serie"`
}

// SensorDataRequest represents the incoming request for sensor data
type SensorDataRequest struct {
	NumeroSerie        string      `json:"numeroSerie"`
	Sensor             string      `json:"sensor"`
	FechaActivacion    string      `json:"fecha_activacion"`
	FechaDesactivacion string      `json:"fecha_desactivacion"`
	Estado             interface{} `json:"estado"` 
}