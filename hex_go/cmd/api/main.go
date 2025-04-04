package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"hex_go/internal/application/services"
	"hex_go/internal/infrastructure/controllers"
	"hex_go/internal/infrastructure/persistence"
	"hex_go/pkg/config"
	"hex_go/pkg/rabbitmq"
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

	// Initialize RabbitMQ client
	rabbitClient, err := rabbitmq.NewRabbitMQClient(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
		log.Printf("Continuing without RabbitMQ integration")
		rabbitClient = nil
	} else {
		defer rabbitClient.Close()
	}

	// Initialize service
	sensorService := services.NewSensorService(repository, rabbitClient)

	// Initialize controller
	sensorController := controllers.NewSensorController(sensorService)

	// Set up router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/api/sensors", sensorController.CreateSensorData).Methods("POST")

	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Apply CORS middleware to router
	handler := c.Handler(router)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}