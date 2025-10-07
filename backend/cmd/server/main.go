package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"genshin-quiz/config"
	"genshin-quiz/internal/webserver"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize configuration
	cfg := config.NewApp()

	// Initialize server
	server := webserver.NewServer(cfg)
	// Start server in a goroutine
	server.Start()
}
