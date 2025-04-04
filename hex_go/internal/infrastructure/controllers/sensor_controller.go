package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	
	// Log the received JSON for debugging
	log.Printf("Received JSON: %s", string(body))

	// Decode the request body
	err = json.Unmarshal(body, &sensorData)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log the parsed data
	log.Printf("Parsed sensor data: %+v", sensorData)

	// Process the sensor data
	err = c.sensorService.ProcessSensorData(&sensorData)
	if err != nil {
		log.Printf("Error processing sensor data: %v", err)
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