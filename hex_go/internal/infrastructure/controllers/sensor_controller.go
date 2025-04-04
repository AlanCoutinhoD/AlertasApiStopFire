package controllers

import (
	"encoding/json"
	"net/http"

	"hex_go/internal/application/services"
	"hex_go/internal/domain/entities"
)

// SensorController handles HTTP requests for sensors
type SensorController struct {
	sensorService *services.SensorService
}

// NewSensorController creates a new sensor controller
func NewSensorController(sensorService *services.SensorService) *SensorController {
	return &SensorController{
		sensorService: sensorService,
	}
}

// CreateSensorData handles the creation of sensor data
func (c *SensorController) CreateSensorData(w http.ResponseWriter, r *http.Request) {
	var sensorData entities.SensorDataRequest

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&sensorData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process the sensor data
	err = c.sensorService.ProcessSensorData(&sensorData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sensor data created successfully",
	})
}