package main

import (
	"github.com/serediukit/civix-backend/internal/server"
	"log"

	"github.com/joho/godotenv"
	"github.com/serediukit/civix-backend/internal/config"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	s := server.NewServer(cfg)

	if err = s.Run(); err != nil {
		log.Fatalf("Failed server: %v", err)
	}

	log.Println("Server exiting")
}
