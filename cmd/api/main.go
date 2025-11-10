package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/serediukit/civix-backend/internal/config"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Initialize configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server := NewServer(config)

	server.setup()

	if err = server.run(); err != nil {
		log.Fatalf("Failed server: %v", err)
	}

	log.Println("Server exiting")
}
