package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"hex_go/internal/application/services"
	"hex_go/internal/infrastructure/controllers"
	"hex_go/internal/infrastructure/persistence"
	"hex_go/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repository := persistence.NewMySQLRepository(db)

	// Initialize service
	sensorService := services.NewSensorService(repository)

	// Initialize controller
	sensorController := controllers.NewSensorController(sensorService)

	// Set up router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/api/sensors", sensorController.CreateSensorData).Methods("POST")

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}