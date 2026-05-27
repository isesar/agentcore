package main

import (
	"agentcore/config"
	"agentcore/db"
	"log"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	if err := db.RunMigrations("db/migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully")
}
