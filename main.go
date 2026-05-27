package main

import (
	"agentcore/config"
	"agentcore/db"
	"agentcore/server"
	"log"
)

func main() {
	// Load environment variables
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	// Initialize database connection
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Run database migrations at startup
	if err := db.RunMigrations("db/migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
