package main

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/db"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	if err := db.InitDB(cfg); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	log.Println("Database migration completed!")
}
