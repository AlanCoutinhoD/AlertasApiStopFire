package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"hex_go/internal/domain/entities"
	"hex_go/internal/domain/ports"
)

type SensorController struct {
	sensorService ports.SensorServicePort
}

func NewSensorController(sensorService ports.SensorServicePort) *SensorController {
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

// GetUserAlerts handles retrieving all alerts for a user
func (c *SensorController) GetUserAlerts(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameter
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}
	
	// Convert user ID to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}
	
	// Get alerts for this user
	alerts, err := c.sensorService.GetUserAlerts(userID)
	if err != nil {
		log.Printf("Error getting user alerts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return alerts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}