package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
    // Try to load .env file from multiple possible locations
    envPaths := []string{
        ".env",
        "../.env",
        "../../.env",
    
    }
    
    var loadErr error
    for _, path := range envPaths {
        err := godotenv.Load(path)
        if err == nil {
            log.Printf("Successfully loaded .env from %s", path)
            loadErr = nil
            break
        }
        loadErr = err
    }
    
    if loadErr != nil {
        log.Printf("Warning: Error loading .env file: %v", loadErr)
        log.Printf("Using default or environment values instead")
    }

    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBUser:     getEnv("DB_USER", "root"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "stopfire"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
    }
}

// ConnectDB establishes a connection to the database
func (c *Config) ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	
	log.Printf("Connecting to database: %s on %s:%s", c.DBName, c.DBHost, c.DBPort)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	
	if err = db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}